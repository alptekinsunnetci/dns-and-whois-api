package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/miekg/dns"
	"golang.org/x/net/idna"
)

var (
	timeoutMs          int
	dnsServer          string
	cache              = make(map[string][]dns.RR)
	cacheTTL           = 5 * time.Minute
	cacheLock          sync.RWMutex
	requestCount       = make(map[string]int)
	rateLimit          = 100
	rateLimitDuration  = time.Minute
	logFile            *os.File
	allowedRecordTypes = map[string]bool{
		"A":     true,
		"AAAA":  true,
		"SOA":   true,
		"MX":    true,
		"NS":    true,
		"CNAME": true,
		"TXT":   true,
	}
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Question struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Class string `json:"class"`
}

type Section struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Class    string `json:"class"`
	Ttl      uint32 `json:"ttl"`
	Rdlength uint16 `json:"rdlength"`
	Rdata    string `json:"rdata"`
}

type Message struct {
	Question   []*Question `json:"question"`
	Answer     []*Section  `json:"answer"`
	Authority  []*Section  `json:"authority,omitempty"`
	Additional []*Section  `json:"additional,omitempty"`
}

type WhoisResult struct {
	Server string `json:"server"`
	Data   string `json:"data"`
}

var domainRegex = regexp.MustCompile(`^[a-zA-Z0-9.-]+$`)

func isValidDomain(domain string) bool {
	return domainRegex.MatchString(domain)
}

func rdata(RR dns.RR) string {
	return strings.Replace(RR.String(), RR.Header().String(), "", -1)
}

func errorResponse(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "dns-api")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Error{Code: code, Message: message})
	logError(fmt.Errorf("Error %d: %s", code, message))
}

func logError(err error) {
	if logFile != nil {
		log.SetOutput(logFile)
		log.Println(err)
	}
}

func resolve(w http.ResponseWriter, r *http.Request, domain string, querytype uint16) ([]dns.RR, error) {
	cacheKey := fmt.Sprintf("%s-%d", domain, querytype)

	cacheLock.RLock()
	if answers, found := cache[cacheKey]; found {
		cacheLock.RUnlock()
		return answers, nil
	}
	cacheLock.RUnlock()

	m := new(dns.Msg)
	m.Id = dns.Id()
	m.RecursionDesired = true
	m.Question = make([]dns.Question, 1)
	m.Question[0] = dns.Question{domain, querytype, dns.ClassINET}

	c := new(dns.Client)
	c.Dialer = &net.Dialer{
		Timeout: time.Duration(timeoutMs) * time.Millisecond,
	}

Redo:
	if in, _, err := c.Exchange(m, dnsServer); err == nil {
		if in.MsgHdr.Truncated {
			c.Net = "tcp"
			goto Redo
		}

		if in.MsgHdr.Rcode != dns.RcodeSuccess {
			return nil, fmt.Errorf("DNS query failed with Rcode: %d", in.MsgHdr.Rcode)
		}

		cacheLock.Lock()
		cache[cacheKey] = in.Answer
		cacheLock.Unlock()

		time.AfterFunc(cacheTTL, func() {
			cacheLock.Lock()
			delete(cache, cacheKey)
			cacheLock.Unlock()
		})

		return in.Answer, nil
	} else {
		return nil, fmt.Errorf("DNS server could not be reached")
	}
}

func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr

		requestCount[clientIP]++
		if requestCount[clientIP] > rateLimit {
			errorResponse(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return
		}

		go func() {
			time.Sleep(rateLimitDuration)
			requestCount[clientIP] = 0
		}()

		next.ServeHTTP(w, r)
	})
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	_, err := resolve(w, r, "iana.org", dns.TypeA)
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, "DNS server is not reachable")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func whoisHandler(w http.ResponseWriter, r *http.Request) {
	domain := chi.URLParam(r, "domain")

	if !isValidDomain(domain) {
		errorResponse(w, http.StatusBadRequest, "Input string is not a well-formed domain name")
		return
	}

	tld := strings.TrimPrefix(domain, "www.")
	tld = strings.Split(tld, ".")[1]

	whoisServer := getWhoisServerFromIANA(tld)
	if whoisServer == "" {
		errorResponse(w, http.StatusNotFound, "WHOIS server not found")
		return
	}

	whoisResults := queryWhoisServer(domain, whoisServer)
	if whoisResults == nil {
		errorResponse(w, http.StatusInternalServerError, "Failed to retrieve WHOIS data")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(whoisResults)
}

func getWhoisServerFromIANA(tld string) string {
	conn, err := net.Dial("tcp", "whois.iana.org:43")
	if err != nil {
		log.Println("Could not connect to WHOIS server:", err)
		return ""
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\r\n", tld)
	if err != nil {
		log.Println("Could not write to WHOIS server:", err)
		return ""
	}

	var response bytes.Buffer
	io.Copy(&response, conn)

	lines := strings.Split(response.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "whois:") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func queryWhoisServer(domain, whoisServer string) []WhoisResult {
	var results []WhoisResult

	for {
		conn, err := net.Dial("tcp", whoisServer+":43")
		if err != nil {
			log.Println("Could not connect to WHOIS server:", err)
			return nil
		}
		defer conn.Close()

		_, err = fmt.Fprintf(conn, "%s\r\n", domain)
		if err != nil {
			log.Println("Could not write to WHOIS server:", err)
			return nil
		}

		var response bytes.Buffer
		io.Copy(&response, conn)

		results = append(results, WhoisResult{Server: whoisServer, Data: response.String()})

		lines := strings.Split(response.String(), "\n")
		nextWhoisServer := ""
		for _, line := range lines {
			if strings.Contains(line, "Registrar WHOIS Server:") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					nextWhoisServer = strings.TrimSpace(parts[1])
					break
				}
			}
		}

		if nextWhoisServer == "" {
			break
		}

		whoisServer = nextWhoisServer
	}

	return results
}

func multiQueryHandler(w http.ResponseWriter, r *http.Request) {
	domain := dns.Fqdn(chi.URLParam(r, "domain"))

	if !isValidDomain(domain) {
		errorResponse(w, http.StatusBadRequest, "Input string is not a well-formed domain name")
		return
	}

	if domain, err := idna.ToASCII(domain); err == nil {
		if _, isDomain := dns.IsDomainName(domain); isDomain {
			recordTypes := []uint16{dns.TypeSOA, dns.TypeA, dns.TypeCNAME, dns.TypeTXT, dns.TypeMX, dns.TypeNS}
			var allAnswers []dns.RR

			for _, recordType := range recordTypes {
				answers, err := resolve(w, r, domain, recordType)
				if err == nil {
					allAnswers = append(allAnswers, answers...)
				}
			}

			jsonify(w, r, []dns.Question{{Name: domain}}, allAnswers, nil, nil)
		} else {
			errorResponse(w, http.StatusBadRequest, "Input string is not a well-formed domain name")
		}
	} else {
		errorResponse(w, http.StatusBadRequest, "Input string could not be parsed")
	}
}

func queryHandler(w http.ResponseWriter, r *http.Request) {
	domain := dns.Fqdn(chi.URLParam(r, "domain"))
	querytype := chi.URLParam(r, "querytype")

	if !allowedRecordTypes[strings.ToUpper(querytype)] {
		errorResponse(w, http.StatusForbidden, "This record type is not allowed")
		return
	}

	if !isValidDomain(domain) {
		errorResponse(w, http.StatusBadRequest, "Input string is not a well-formed domain name")
		return
	}

	if domain, err := idna.ToASCII(domain); err == nil {
		if _, isDomain := dns.IsDomainName(domain); isDomain {
			if querytype, ok := dns.StringToType[strings.ToUpper(querytype)]; ok {
				answers, err := resolve(w, r, domain, querytype)
				if err == nil {
					jsonify(w, r, []dns.Question{{Name: domain, Qtype: querytype}}, answers, nil, nil)
				} else {
					errorResponse(w, http.StatusInternalServerError, err.Error())
				}
			} else {
				errorResponse(w, http.StatusNotFound, "Invalid DNS query type")
			}
		} else {
			errorResponse(w, http.StatusBadRequest, "Input string is not a well-formed domain name")
		}
	} else {
		errorResponse(w, http.StatusBadRequest, "Input string could not be parsed")
	}
}

func ptrHandler(w http.ResponseWriter, r *http.Request) {
	ip := chi.URLParam(r, "ip")

	if arpa, err := dns.ReverseAddr(ip); err == nil {
		answers, err := resolve(w, r, arpa, dns.TypePTR)
		if err == nil {
			jsonify(w, r, []dns.Question{{Name: arpa, Qtype: dns.TypePTR}}, answers, nil, nil)
		} else {
			errorResponse(w, http.StatusInternalServerError, err.Error())
		}
	} else {
		errorResponse(w, http.StatusBadRequest, "Input string is not a valid IP address")
	}
}

func jsonify(w http.ResponseWriter, r *http.Request, question []dns.Question, answer []dns.RR, authority []dns.RR, additional []dns.RR) {
	var answerArray, authorityArray, additionalArray []*Section

	callback := r.URL.Query().Get("callback")

	for _, answer := range answer {
		answerArray = append(answerArray, &Section{answer.Header().Name, dns.TypeToString[answer.Header().Rrtype], dns.ClassToString[answer.Header().Class], answer.Header().Ttl, answer.Header().Rdlength, rdata(answer)})
	}

	for _, authority := range authority {
		authorityArray = append(authorityArray, &Section{authority.Header().Name, dns.TypeToString[authority.Header().Rrtype], dns.ClassToString[authority.Header().Class], authority.Header().Ttl, authority.Header().Rdlength, rdata(authority)})
	}

	for _, additional := range additional {
		additionalArray = append(additionalArray, &Section{additional.Header().Name, dns.TypeToString[additional.Header().Rrtype], dns.ClassToString[additional.Header().Class], additional.Header().Ttl, additional.Header().Rdlength, rdata(additional)})
	}

	if json, err := json.MarshalIndent(Message{[]*Question{{question[0].Name, dns.TypeToString[question[0].Qtype], dns.ClassToString[question[0].Qclass]}}, answerArray, authorityArray, additionalArray}, "", "    "); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Server", "dns-api")
		if callback != "" {
			io.WriteString(w, callback+"("+string(json)+");")
		} else {
			io.WriteString(w, string(json))
		}
	}
}

func main() {
	var err error
	logFile, err = os.OpenFile("dns_api.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Could not open log file")
	}
	defer logFile.Close()

	host := flag.String("host", "127.0.0.1", "Set the server host")
	port := flag.String("port", "8080", "Set the server port")
	flag.IntVar(&timeoutMs, "timeout", 2000, "Set the query timeout in ms")
	flag.StringVar(&dnsServer, "dns-server", "1.1.1.1:53", "Set the DNS server address")
	version := flag.Bool("version", false, "Display version")

	flag.Usage = func() {
		fmt.Println("\nUSAGE:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *version {
		fmt.Println("dns-and-whois-api 1.0.0")
		os.Exit(0)
	}

	address := *host + ":" + *port

	r := chi.NewRouter()
	r.Use(rateLimiter)
	r.Get("/health", healthCheckHandler)
	r.Get("/{domain}/{querytype}", queryHandler)
	r.Get("/{domain}", multiQueryHandler)
	r.Get("/whois/{domain}", whoisHandler)
	r.Get("/ptr/{ip}", ptrHandler)

	fmt.Printf("Listening on: http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, r))
}
