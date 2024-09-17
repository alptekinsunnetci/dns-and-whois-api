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
- `GET /domain/{domain}`: Perform multiple DNS queries (A, AAAA, SOA, MX, NS, CNAME, TXT) for the specified domain.
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
