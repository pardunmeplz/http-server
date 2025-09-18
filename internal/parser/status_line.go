package parser

import (
	"bytes"
	"fmt"
	req "tcpServer/internal/request"
)

const SUB_SEPERATOR = ' '

type ParseVerbState struct{}

func (verbParser *ParseVerbState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.IndexByte(data, SUB_SEPERATOR)
	if bytes.Contains(data, []byte{'\r', '\n'}) && sepIndex == -1 {
		parser.state = &ErrorState{MALFORMED_REQUEST_LINE_ERR}
		return 0, fmt.Errorf("ERR_VERB_STATE: %s", MALFORMED_REQUEST_LINE_ERR)
	}

	if sepIndex == -1 {
		return 0, nil
	}

	verb := data[:sepIndex]
	parser.Request.RequestLine = req.RequestLine{HttpVersion: "", RequestTarget: "", Method: string(verb)}
	parser.state = &ParseTargetState{}
	return len(verb), nil
}

type ParseTargetState struct{}

func (targetParser *ParseTargetState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.IndexByte(data, SUB_SEPERATOR)
	if bytes.Contains(data, []byte(SEPERATOR)) && sepIndex == -1 {
		parser.state = &ErrorState{MALFORMED_REQUEST_LINE_ERR}
		return 0, fmt.Errorf("ERR_TARGET_STATE: %s", MALFORMED_REQUEST_LINE_ERR)
	}

	if sepIndex == -1 {
		return 0, nil
	}

	location := data[:sepIndex]
	parser.Request.RequestLine.RequestTarget = string(location)
	parser.state = &ParseVerbState{}
	return len(location), nil
}

type ParseVersionState struct{}

func (locationParser *ParseVersionState) parse(data []byte, parser *Parser) (int, error) {
	sepIndex := bytes.Index(data, []byte(SEPERATOR))
	if sepIndex == -1 {
		return 0, nil
	}

	version := data[:sepIndex]
	parser.Request.RequestLine.HttpVersion = string(version)
	parser.state = &DoneState{}
	return len(version), nil
}
