package request

import "strings"

type Request struct {
	RequestLine RequestLine
	Headers     Headers
	Body        []byte
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Headers map[string]string

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name string, value string) {
	if _, ok := h[name]; ok {
		h[strings.ToLower(name)] = h[name] + ", " + value
	} else {
		h[strings.ToLower(name)] = value
	}
}

func NewHeaders() Headers {
	return make(Headers)
}
