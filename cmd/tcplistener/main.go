package main

import (
	"fmt"
	"log"
	"net"
	req "tcpServer/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:42069")
	defer listener.Close()
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		request, err := req.RequestFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}
		printRequest(request)

	}

}

func printRequest(request *req.Request) {
	fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s",
		request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
}
