package parser

import (
	"fmt"
	"strconv"
)

type ParseBody struct{}

func (bodyParser *ParseBody) parse(data []byte, parser *Parser) (int, error) {
	length, err := strconv.Atoi(parser.Request.Headers.Get("content-length"))
	if err != nil {
		parser.state = &ErrorState{INVALID_CONTENT_LENGTH}
		return 0, fmt.Errorf("%s", INVALID_CONTENT_LENGTH)
	}
	parser.Request.Body = append(parser.Request.Body, data...)
	if len(parser.Request.Body) == length {
		parser.state = &DoneState{}
	} else if len(parser.Request.Body) > length {
		parser.state = &ErrorState{CONTENT_LENGTH_MISMATCH}
		return len(data), fmt.Errorf("%s", CONTENT_LENGTH_MISMATCH)
	} else if len(data) == 0 {
		parser.processing = false
	}
	return len(data), nil
}
