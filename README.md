# TCP HTTP Server (from scratch)

A minimal HTTP/1.1 server and request parser written from scratch in Go. Built while following the Boot.dev HTTP course, this project implements a custom HTTP parser, a lightweight TCP-based HTTP server, and several example handlers (including a basic httpbin proxy and chunked responses).
The project is for exploration and may lack some error handling capabilities / have a minimal subset of the protocol implemented keeping things like keep-alive or trailers out of scope

## Features

- Custom HTTP/1.1 parsing using a state machine
  - Parses request line (method, target, version)
  - Validates header names and parses headers
  - Optional body parsing via Content-Length
- Simple TCP server with concurrent connection handling
- Regex-based router with not-found handler
- Graceful shutdown on SIGINT/SIGTERM
- Response writer with status line, headers, body
- Basic chunked transfer encoding support and streaming proxy to httpbin
- Example routes
  - `/yourproblem` → 400 response
  - `/myproblem` → 500 response
  - `/httpbin/*` → basic proxy to httpbin.org (chunked streaming)
  - default route → 200 response

## Status

- Parser: largely complete and tested
- Server: work-in-progress; functional but intentionally minimal and missing many production features

## Getting Started

### Prerequisites
- Go 1.21+ (module requires 1.24 in `go.mod`; use latest stable Go)

### Build

```bash
# Build the simple TCP listener (prints parsed requests)
go build -o tcplistener ./cmd/tcplistener

# Build the HTTP server
go build -o httpserver ./cmd/httpserver
```

### Run

```bash
# Terminal 1: run the HTTP server (listens on :42069)
./httpserver

# Terminal 2: try some requests
curl -v http://localhost:42069/
curl -v http://localhost:42069/yourproblem
curl -v http://localhost:42069/myproblem
curl -v http://localhost:42069/httpbin/get

# Or run the bare TCP listener and send raw HTTP
./tcplistener
printf 'GET / HTTP/1.1\r\nHost: localhost\r\n\r\n' | nc -N localhost 42069
```

## Tests

Unit tests cover the HTTP parser, including request line parsing, headers, and bodies with varying chunk sizes.

```bash
go test ./internal/parser/... -v
```
## Project Structure

```
cmd/
  httpserver/     # runnable HTTP server with example routes
  tcplistener/    # minimal TCP listener that prints parsed request
internal/
  parser/         # HTTP request parser (state machine + tests)
  request/        # request types and helpers
  response/       # response writer and status handling
  server/         # server loop, router, and wiring
```

## Inspiration / Credit

Created while following the Boot.dev HTTP course. This repo is intentionally educational, focusing on understanding HTTP over TCP and parser design rather than providing a production web server.

## Next Steps / Ideas

- Improve chunked responses and add chunked request parsing
- Middleware-style handler composition
- Basic router with path params
- Connection timeouts and keep-alive support
