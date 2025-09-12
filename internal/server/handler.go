package server

import (
	"tcpServer/internal/request"
	"tcpServer/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}
type Handler func(w response.Writer, req *request.Request)

func (herr *HandlerError) writeError(w response.Writer) {
	headers := response.GetDefaultHeaders(len(herr.Message))
	w.WriteStatusLine(herr.Code)
	w.WriteHeaders(headers)
	w.WriteBody([]byte(herr.Message))
}
