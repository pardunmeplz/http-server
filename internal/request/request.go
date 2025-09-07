package request

import (
	"fmt"
	"io"
	"strings"
	h "tcpServer/internal/headers"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	Headers     h.Headers
	Body        []byte
	State       parserState
}

type parserState int

const (
	INITIALIZED parserState = 0
	HEADERS     parserState = 1
	BODY        parserState = 2
	DONE        parserState = 3
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const SEPERATOR = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := []byte{}
	request := Request{RequestLine{}, make(h.Headers), nil, INITIALIZED}
	totConsumed := 0
	for request.State != DONE {
		chunk := make([]byte, 8)
		size, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, chunk[:size]...)
		consumed, err := request.parse(buffer[totConsumed:])
		if err != nil {
			return nil, err
		}
		totConsumed += consumed
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totConsumed := 0
	if r.State == INITIALIZED {
		line, consumed, error := ParseRequestLine(string(data))
		totConsumed += consumed
		if error != nil {
			return totConsumed, error
		}
		if line != nil {
			r.RequestLine = *line
			r.State = HEADERS
		}
	}
	if r.State == HEADERS {
		for {
			consumed, done, error := r.Headers.Parse(data[totConsumed:])
			totConsumed += consumed
			if error != nil {
				return totConsumed, error
			}
			if done {
				r.State = DONE
				break
			}
			if consumed == 0 {
				break
			}
			fmt.Println(consumed)
		}
	}
	return totConsumed, nil
}

func ParseRequestLine(rawStr string) (*RequestLine, int, error) {

	lineEnd := strings.Index(rawStr, SEPERATOR)
	// not enough data to parse request line
	if lineEnd == -1 {
		return nil, 0, nil
	}
	line := rawStr[:lineEnd]

	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, lineEnd + len(SEPERATOR), fmt.Errorf("Invalid request line %d", lineEnd)
	}
	subParts := strings.Split(parts[2], "/")
	if len(subParts) != 2 || subParts[0] != "HTTP" {
		return nil, lineEnd + len(SEPERATOR), fmt.Errorf("unexpected HTTP version %s", parts[2])
	}
	for _, r := range parts[0] {
		if !unicode.IsUpper(r) {
			return nil, lineEnd + len(SEPERATOR), fmt.Errorf("unexpected Method %s", parts[0])
		}
	}
	return &RequestLine{subParts[1], parts[1], parts[0]}, lineEnd + len(SEPERATOR), nil
}
