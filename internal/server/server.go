package server

import (
	"fmt"
	"log"
	"net"
	"tcpServer/internal/parser"
	"tcpServer/internal/response"
)

type Server struct {
	listener net.Listener
	handler  Handler
	parser   parser.Parser
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener, handler, parser.Parser{}}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		defer conn.Close()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := s.parser.ParseFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}

	writer := response.Writer{Writer: conn}
	s.handler(writer, req)

}
