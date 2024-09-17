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
