package parser

import (
	"bytes"
	"fmt"
	"strings"
	req "tcpServer/internal/request"
	"unicode"
)

const STATUS_LINE_SEPERATOR = ' '

type ParseVerbState struct{}

func (verbParser *ParseVerbState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.IndexByte(data, STATUS_LINE_SEPERATOR)
	if bytes.Contains(data, []byte{'\r', '\n'}) && sepIndex == -1 {
		parser.state = &ErrorState{MALFORMED_REQUEST_LINE_ERR}
		return 0, fmt.Errorf("ERR_VERB_STATE: %s", MALFORMED_REQUEST_LINE_ERR)
	}

	if sepIndex == -1 {
		parser.processing = false
		return 0, nil
	}

	verb := data[:sepIndex]

	for _, b := range verb {
		if !unicode.IsUpper(rune(b)) {
			parser.state = &ErrorState{INVALID_VERB_ERR}
			return 0, fmt.Errorf("ERR_VERB_STATE: %s", INVALID_VERB_ERR)
		}
	}

	parser.Request.RequestLine = req.RequestLine{HttpVersion: "", RequestTarget: "", Method: string(verb)}
	parser.state = &ParseTargetState{}
	return len(verb) + 1, nil
}

type ParseTargetState struct{}

func (targetParser *ParseTargetState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.IndexByte(data, STATUS_LINE_SEPERATOR)
	if bytes.Contains(data, []byte(SEPERATOR)) && sepIndex == -1 {
		parser.state = &ErrorState{MALFORMED_REQUEST_LINE_ERR}
		return 0, fmt.Errorf("ERR_TARGET_STATE: %s", MALFORMED_REQUEST_LINE_ERR)
	}

	if sepIndex == -1 {
		parser.processing = false
		return 0, nil
	}

	location := data[:sepIndex]
	parser.Request.RequestLine.RequestTarget = string(location)
	parser.state = &ParseVersionState{}
	return len(location) + 1, nil
}

type ParseVersionState struct{}

func (versionParser *ParseVersionState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.Index(data, []byte(SEPERATOR))
	if sepIndex == -1 {
		parser.processing = false
		return 0, nil
	}

	version := string(data[:sepIndex])
	if !strings.HasPrefix(version, "HTTP/") {
		parser.state = &ErrorState{UNEXPECTED_VERSION_ERR}
		return 0, fmt.Errorf("ERR_VERSION_STATE: %s", UNEXPECTED_VERSION_ERR)
	}
	parser.Request.RequestLine.HttpVersion = version[5:]
	parser.state = &ParseHeaders{}
	return len(version) + len(SEPERATOR), nil
}
