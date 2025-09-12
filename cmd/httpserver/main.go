package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tcpServer/internal/request"
	"tcpServer/internal/response"
	sv "tcpServer/internal/server"
)

const port = 42069

func main() {
	server, err := sv.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w response.Writer, req *request.Request) {
	message := ""
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		message = `<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`
		w.WriteStatusLine(response.BAD_REQUEST)
	case "/myproblem":
		message = `<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`
		w.WriteStatusLine(response.INTERNAL_SERVER_ERROR)
	default:
		message = `<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`
		w.WriteStatusLine(response.OK)
	}
	headers := response.GetDefaultHeaders(len(message))
	headers.Set("content-type", "text/html")
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))

}
