# DNS & WHOIS API

This project is a Go-based API server for handling DNS and WHOIS queries. It provides endpoints for various DNS record types (A, AAAA, SOA, MX, NS, CNAME, TXT, and PTR) and WHOIS lookups, supporting caching and rate-limiting.

## Features

- **DNS Querying**: Resolve multiple DNS record types, including A, AAAA, SOA, MX, NS, CNAME, TXT, and PTR records.
- **WHOIS Lookup**: Query WHOIS servers to retrieve domain ownership information, following any necessary server redirections.
- **Caching**: DNS responses are cached for a configurable period to reduce load on the upstream DNS servers.
- **Rate Limiting**: Requests from the same IP address are rate-limited to prevent abuse.
- **Health Check**: Basic health check endpoint that queries a DNS server to verify availability.

## Endpoints

### DNS Endpoints

- `GET /{domain}/{querytype}`: Perform a DNS query for the specified domain and query type.
- `GET /{domain}`: Perform multiple DNS queries (A, AAAA, SOA, MX, NS, CNAME, TXT) for the specified domain.
- `GET /ptr/{ip}`: Perform a reverse DNS (PTR) query for the specified IP address.
- `GET /health`: Health check endpoint that verifies DNS server availability.

### WHOIS Endpoints

- `GET /whois/{domain}`: Perform a WHOIS lookup for the specified domain.

## Usage

### Running the API

To run the server, use the following commands:

```bash
go run main.go
```


## Configuration

- Server Host: -host (default: 127.0.0.1)
- Server Port: -port (default: 8080)
- Query Timeout: -timeout in milliseconds (default: 2000)
- DNS Server: -dns-server address (default: 1.1.1.1:53)


## Sample

```curl
http://127.0.0.1:8080/google.com
```

```json
{
    "question": [
        {
            "name": "google.com.",
            "type": "None",
            "class": ""
        }
    ],
    "answer": [
        {
            "name": "google.com.",
            "type": "SOA",
            "class": "IN",
            "ttl": 60,
            "rdlength": 38,
            "rdata": "ns1.google.com. dns-admin.google.com. 675491404 900 900 1800 60"
        },
        {
            "name": "google.com.",
            "type": "A",
            "class": "IN",
            "ttl": 82,
            "rdlength": 4,
            "rdata": "142.250.187.174"
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 60,
            "rdata": "\"facebook-domain-verification=22rm551cu4k0ab0bxsw536tlds4h95\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 69,
            "rdata": "\"google-site-verification=wD8N7i1JTNTkezJ49swvWW48f8_9xveREV4oB-0Hf5o\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 69,
            "rdata": "\"google-site-verification=4ibFUgB-wXLQ_S7vsXVomSTVamuOXBiVAzpR5IZ87D0\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 43,
            "rdata": "\"apple-domain-verification=30afIBcvSuDV2PLX\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 46,
            "rdata": "\"docusign=1b0a6754-49b1-4db5-8540-d2c12664b289\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 65,
            "rdata": "\"globalsign-smime-dv=CDYX+XFHUw2wml6/Gb8+59BsH31KzUr6c1l2BPvqKX8=\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 69,
            "rdata": "\"google-site-verification=TV9-DBe4R80X4v0M4U_bd_J9cpOJM0nikft0jAgjmsQ\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 36,
            "rdata": "\"v=spf1 include:_spf.google.com ~all\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 46,
            "rdata": "\"docusign=05958488-4752-4ef2-95eb-aa7ba8a3bd0e\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 44,
            "rdata": "\"MS=E4A68B9AB2BB9670BCE15412F62916164C0B20BB\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 94,
            "rdata": "\"cisco-ci-domain-verification=479146de172eb01ddee38b1a455ab9e8bb51542ddd7f1fa298557dfa7b22d963\""
        },
        {
            "name": "google.com.",
            "type": "TXT",
            "class": "IN",
            "ttl": 3587,
            "rdlength": 62,
            "rdata": "\"onetrust-domain-verification=de01ed21f2fa4d8781cbc3ffb89cf4ef\""
        },
        {
            "name": "google.com.",
            "type": "MX",
            "class": "IN",
            "ttl": 300,
            "rdlength": 9,
            "rdata": "10 smtp.google.com."
        },
        {
            "name": "google.com.",
            "type": "NS",
            "class": "IN",
            "ttl": 337856,
            "rdlength": 6,
            "rdata": "ns4.google.com."
        },
        {
            "name": "google.com.",
            "type": "NS",
            "class": "IN",
            "ttl": 337856,
            "rdlength": 6,
            "rdata": "ns2.google.com."
        },
        {
            "name": "google.com.",
            "type": "NS",
            "class": "IN",
            "ttl": 337856,
            "rdlength": 6,
            "rdata": "ns1.google.com."
        },
        {
            "name": "google.com.",
            "type": "NS",
            "class": "IN",
            "ttl": 337856,
            "rdlength": 6,
            "rdata": "ns3.google.com."
        }
    ]
}

```



```curl
http://127.0.0.1:8080/google.com/mx
```

```json
{
    "question": [
        {
            "name": "google.com.",
            "type": "MX",
            "class": ""
        }
    ],
    "answer": [
        {
            "name": "google.com.",
            "type": "MX",
            "class": "IN",
            "ttl": 300,
            "rdlength": 9,
            "rdata": "10 smtp.google.com."
        }
    ]
}
```



```curl
http://127.0.0.1:8080/ptr/1.1.1.1
```

```json
{
    "question": [
        {
            "name": "1.1.1.1.in-addr.arpa.",
            "type": "PTR",
            "class": ""
        }
    ],
    "answer": [
        {
            "name": "1.1.1.1.in-addr.arpa.",
            "type": "PTR",
            "class": "IN",
            "ttl": 321,
            "rdlength": 17,
            "rdata": "one.one.one.one."
        }
    ]
}
```




```curl
http://127.0.0.1:8080/whois/google.com
```

```json
[
  {
    "server": "whois.verisign-grs.com",
    "data": "   Domain Name: GOOGLE.COM\r\n   Registry Domain ID: 2138514_DOMAIN_COM-VRSN\r\n   Registrar WHOIS Server: whois.markmonitor.com\r\n   Registrar URL: http://www.markmonitor.com\r\n   Updated Date: 2019-09-09T15:39:04Z\r\n   Creation Date: 1997-09-15T04:00:00Z\r\n   Registry Expiry Date: 2028-09-14T04:00:00Z\r\n   Registrar: MarkMonitor Inc.\r\n   Registrar IANA ID: 292\r\n   Registrar Abuse Contact Email: abusecomplaints@markmonitor.com\r\n   Registrar Abuse Contact Phone: +1.2086851750\r\n   Domain Status: clientDeleteProhibited https://icann.org/epp#clientDeleteProhibited\r\n   Domain Status: clientTransferProhibited https://icann.org/epp#clientTransferProhibited\r\n   Domain Status: clientUpdateProhibited https://icann.org/epp#clientUpdateProhibited\r\n   Domain Status: serverDeleteProhibited https://icann.org/epp#serverDeleteProhibited\r\n   Domain Status: serverTransferProhibited https://icann.org/epp#serverTransferProhibited\r\n   Domain Status: serverUpdateProhibited https://icann.org/epp#serverUpdateProhibited\r\n   Name Server: NS1.GOOGLE.COM\r\n   Name Server: NS2.GOOGLE.COM\r\n   Name Server: NS3.GOOGLE.COM\r\n   Name Server: NS4.GOOGLE.COM\r\n   DNSSEC: unsigned\r\n   URL of the ICANN Whois Inaccuracy Complaint Form: https://www.icann.org/wicf/\r\n>>> Last update of whois database: 2024-09-17T21:00:22Z <<<\r\n\r\nFor more information on Whois status codes, please visit https://icann.org/epp\r\n\r\nNOTICE: The expiration date displayed in this record is the date the\r\nregistrar's sponsorship of the domain name registration in the registry is\r\ncurrently set to expire. This date does not necessarily reflect the expiration\r\ndate of the domain name registrant's agreement with the sponsoring\r\nregistrar.  Users may consult the sponsoring registrar's Whois database to\r\nview the registrar's reported date of expiration for this registration.\r\n\r\nTERMS OF USE: You are not authorized to access or query our Whois\r\ndatabase through the use of electronic processes that are high-volume and\r\nautomated except as reasonably necessary to register domain names or\r\nmodify existing registrations; the Data in VeriSign Global Registry\r\nServices' (\"VeriSign\") Whois database is provided by VeriSign for\r\ninformation purposes only, and to assist persons in obtaining information\r\nabout or related to a domain name registration record. VeriSign does not\r\nguarantee its accuracy. By submitting a Whois query, you agree to abide\r\nby the following terms of use: You agree that you may use this Data only\r\nfor lawful purposes and that under no circumstances will you use this Data\r\nto: (1) allow, enable, or otherwise support the transmission of mass\r\nunsolicited, commercial advertising or solicitations via e-mail, telephone,\r\nor facsimile; or (2) enable high volume, automated, electronic processes\r\nthat apply to VeriSign (or its computer systems). The compilation,\r\nrepackaging, dissemination or other use of this Data is expressly\r\nprohibited without the prior written consent of VeriSign. You agree not to\r\nuse electronic processes that are automated and high-volume to access or\r\nquery the Whois database except as reasonably necessary to register\r\ndomain names or modify existing registrations. VeriSign reserves the right\r\nto restrict your access to the Whois database in its sole discretion to ensure\r\noperational stability.  VeriSign may restrict or terminate your access to the\r\nWhois database for failure to abide by these terms of use. VeriSign\r\nreserves the right to modify these terms at any time.\r\n\r\nThe Registry database contains ONLY .COM, .NET, .EDU domains and\r\nRegistrars.\r\n"
  },
  {
    "server": "whois.markmonitor.com",
    "data": "Domain Name: google.com\nRegistry Domain ID: 2138514_DOMAIN_COM-VRSN\nRegistrar WHOIS Server: whois.markmonitor.com\nRegistrar URL: http://www.markmonitor.com\nUpdated Date: 2024-08-02T02:17:33+0000\nCreation Date: 1997-09-15T07:00:00+0000\nRegistrar Registration Expiration Date: 2028-09-13T07:00:00+0000\nRegistrar: MarkMonitor, Inc.\nRegistrar IANA ID: 292\nRegistrar Abuse Contact Email: abusecomplaints@markmonitor.com\nRegistrar Abuse Contact Phone: +1.2086851750\nDomain Status: clientUpdateProhibited (https://www.icann.org/epp#clientUpdateProhibited)\nDomain Status: clientTransferProhibited (https://www.icann.org/epp#clientTransferProhibited)\nDomain Status: clientDeleteProhibited (https://www.icann.org/epp#clientDeleteProhibited)\nDomain Status: serverUpdateProhibited (https://www.icann.org/epp#serverUpdateProhibited)\nDomain Status: serverTransferProhibited (https://www.icann.org/epp#serverTransferProhibited)\nDomain Status: serverDeleteProhibited (https://www.icann.org/epp#serverDeleteProhibited)\nRegistrant Organization: Google LLC\nRegistrant State/Province: CA\nRegistrant Country: US\nRegistrant Email: Select Request Email Form at https://domains.markmonitor.com/whois/google.com\nAdmin Organization: Google LLC\nAdmin State/Province: CA\nAdmin Country: US\nAdmin Email: Select Request Email Form at https://domains.markmonitor.com/whois/google.com\nTech Organization: Google LLC\nTech State/Province: CA\nTech Country: US\nTech Email: Select Request Email Form at https://domains.markmonitor.com/whois/google.com\nName Server: ns2.google.com\nName Server: ns3.google.com\nName Server: ns4.google.com\nName Server: ns1.google.com\nDNSSEC: unsigned\nURL of the ICANN WHOIS Data Problem Reporting System: http://wdprs.internic.net/\n>>> Last update of WHOIS database: 2024-09-17T20:58:20+0000 <<<\n\nFor more information on WHOIS status codes, please visit:\n  https://www.icann.org/resources/pages/epp-status-codes\n\nIf you wish to contact this domain’s Registrant, Administrative, or Technical\ncontact, and such email address is not visible above, you may do so via our web\nform, pursuant to ICANN’s Temporary Specification. To verify that you are not a\nrobot, please enter your email address to receive a link to a page that\nfacilitates email communication with the relevant contact(s).\n\nWeb-based WHOIS:\n  https://domains.markmonitor.com/whois\n\nIf you have a legitimate interest in viewing the non-public WHOIS details, send\nyour request and the reasons for your request to whoisrequest@markmonitor.com\nand specify the domain name in the subject line. We will review that request and\nmay ask for supporting documentation and explanation.\n\nThe data in MarkMonitor’s WHOIS database is provided for information purposes,\nand to assist persons in obtaining information about or related to a domain\nname’s registration record. While MarkMonitor believes the data to be accurate,\nthe data is provided \"as is\" with no guarantee or warranties regarding its\naccuracy.\n\nBy submitting a WHOIS query, you agree that you will use this data only for\nlawful purposes and that, under no circumstances will you use this data to:\n  (1) allow, enable, or otherwise support the transmission by email, telephone,\nor facsimile of mass, unsolicited, commercial advertising, or spam; or\n  (2) enable high volume, automated, or electronic processes that send queries,\ndata, or email to MarkMonitor (or its systems) or the domain name contacts (or\nits systems).\n\nMarkMonitor reserves the right to modify these terms at any time.\n\nBy submitting this query, you agree to abide by this policy.\n\nMarkMonitor Domain Management(TM)\nProtecting companies and consumers in a digital world.\n\nVisit MarkMonitor at https://www.markmonitor.com\nContact us at +1.8007459229\nIn Europe, at +44.02032062220\n--\n"
  },
  {
    "server": "whois.markmonitor.com",
    "data": ""
  }
]
```

