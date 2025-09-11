package server

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"tcpServer/internal/request"
	"tcpServer/internal/response"
)

type Server struct {
	listener net.Listener
}

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener}
	go server.listen(handler)

	return server, nil
}

func (s *Server) Close() error {
	return s.listener.Close()
}
func (s *Server) listen(handler Handler) {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		s.handle(conn, handler)
	}
}

func (s *Server) handle(conn net.Conn, handler Handler) {
	req, err := request.RequestFromReader(conn)
	if err != nil {
		log.Fatal(err)
	}

	buffer := bytes.Buffer{}
	herr := handler(&buffer, req)
	if herr != nil {
		herr.writeError(conn)
	} else {
		headers := response.GetDefaultHeaders(buffer.Len())
		response.WriteStatusLine(conn, response.OK)
		response.WriteHeaders(conn, headers)
		conn.Write(buffer.Bytes())
	}

}
