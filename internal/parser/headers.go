package parser

import (
	"bytes"
	"fmt"
	"strings"
	req "tcpServer/internal/request"
	"unicode"
)

type ParseHeaders struct{}

const HEADER_SEPERATOR = ':'

func (pareHeaders *ParseHeaders) parse(data []byte, parser *Parser) (int, error) {

	sepIndex := bytes.Index(data, []byte(SEPERATOR))
	if sepIndex == -1 {
		return 0, nil
	}
	if sepIndex == 0 {
		parser.state = &DoneState{}
		return len(SEPERATOR), nil
	}
	if parser.Request.Headers == nil {
		parser.Request.Headers = make(req.Headers)
	}

	nameEnd := bytes.IndexByte(data, HEADER_SEPERATOR)
	name := string(bytes.TrimLeft(data[:nameEnd], " "))
	if !isValidFieldName(name) {
		return 0, fmt.Errorf("%s", INVALID_HEADER_NAME)
	}
	parser.Request.Headers.Set(name, string(bytes.Trim(data[nameEnd+1:sepIndex], " ")))

	return sepIndex + len(SEPERATOR), nil
}

func isValidFieldName(name string) bool {
	if len(name) < 1 {
		return false
	}

	if strings.Index(name, " ") != -1 {
		return false
	}

	for _, ch := range name {
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			if unicode.IsDigit(ch) || unicode.IsLetter(ch) {
				continue
			}
			return false
		}
	}

	return true
}
