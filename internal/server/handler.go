package server

import (
	"io"
	"tcpServer/internal/request"
	"tcpServer/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}
type Handler func(w io.Writer, req *request.Request) *HandlerError

func (herr *HandlerError) writeError(conn io.Writer) {
	headers := response.GetDefaultHeaders(len(herr.Message))
	response.WriteStatusLine(conn, herr.Code)
	response.WriteHeaders(conn, headers)
	conn.Write([]byte(herr.Message))

}
