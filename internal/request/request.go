package request

import (
	"fmt"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
}

type parserState int

const (
	INITIALIZED parserState = 0
	DONE        parserState = 1
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const SEPERATOR = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	buffer := []byte{}
	request := Request{RequestLine{}, INITIALIZED}
	for request.state != DONE {
		chunk := make([]byte, 8)
		size, err := reader.Read(chunk)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, chunk[:size]...)
		consumed, err := request.parse(buffer)
		if err != nil {
			return nil, err
		}
		if consumed != 0 {
			buffer = buffer[:consumed]
		}
	}

	return &request, nil
}

func (r *Request) parse(data []byte) (int, error) {
	line, consumed, error := ParseRequestLine(string(data))
	if error != nil {
		return 0, error
	}
	if consumed == 0 {
		return 0, nil
	}
	r.state = DONE
	r.RequestLine = *line
	return consumed, nil
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
