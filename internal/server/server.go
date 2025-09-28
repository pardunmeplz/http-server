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
	router   *Router
	parser   parser.Parser
}

func Serve(port int, router *Router) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener, router, parser.Parser{}}
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
			log.Print(err)
			writer := response.Writer{Writer: conn}
			error := HandlerError{response.INTERNAL_SERVER_ERROR, err.Error()}
			error.writeError(writer)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	req, err := s.parser.ParseFromReader(conn)
	writer := response.Writer{Writer: conn}
	if err != nil {
		error := HandlerError{response.INTERNAL_SERVER_ERROR, err.Error()}
		error.writeError(writer)
		return
	}
	// performing recover check before running handler
	defer func() {
		if r := recover(); r != nil {
			error := HandlerError{response.INTERNAL_SERVER_ERROR, "something went wrong!"}
			error.writeError(writer)
		}
	}()
	s.router.route(req.RequestLine.RequestTarget)(writer, req)
}
