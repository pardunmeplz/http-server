package server

import (
	"fmt"
	"log"
	"net"
	"tcpServer/internal/response"
)

type Server struct {
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener}
	go server.listen()

	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response.WriteStatusLine(conn, response.OK)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
}
