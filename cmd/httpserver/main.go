package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"tcpServer/internal/request"
	"tcpServer/internal/response"
	sv "tcpServer/internal/server"
)

const port = 42069

func main() {

	router := &sv.Router{}
	router.Register("^/yourproblem$", yourProblem) // equality check
	router.Register("^/myproblem$", myProblem)     // equality check
	router.Register("^/httpbin", httpBin)          // prefix check, you can add anything after /httpbin and it will match
	router.RegisterNotFound(defaultResp)

	server, err := sv.Serve(port, router)
	defer server.Close()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server started on port", port)

	// code to close the server gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

const httpBinLink = "https://httpbin.org"

func httpBin(w response.Writer, req *request.Request) {
	// route
	path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
	w.WriteStatusLine(response.OK)

	// headers
	headers := response.GetDefaultHeaders(0)
	delete(headers, "content-length")
	headers.Set("transfer-encoding", "chunked")
	headers.Set("trailer", "X-Content-Length")
	w.WriteHeaders(headers)

	resp, err := http.Get(httpBinLink + path)
	if err != nil {
		log.Print(err)
		return
	}

	size, err := w.WriteChunkedBody(resp)
	if err != nil {
		log.Print(err)
	}

	trailers := make(request.Headers)
	sizeStr := strconv.Itoa(size)
	trailers.Set("X-Content-Length", sizeStr)

	w.WriteTrailers(trailers)

}

func yourProblem(w response.Writer, req *request.Request) {
	message := `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
	w.WriteStatusLine(response.BAD_REQUEST)
	headers := response.GetDefaultHeaders(len(message))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func myProblem(w response.Writer, req *request.Request) {
	message := `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
	w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)

	headers := response.GetDefaultHeaders(len(message))
	headers.Set("content-type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func defaultResp(w response.Writer, req *request.Request) {
	message := `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
	w.WriteStatusLine(response.OK)
	headers := response.GetDefaultHeaders(len(message))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}
