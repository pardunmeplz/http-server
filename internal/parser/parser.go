package parser

import (
	"fmt"
	"io"
	req "tcpServer/internal/request"
)

const (
	MALFORMED_REQUEST_LINE_ERR = "Malformed request line"
	UNEXPECTED_VERSION_ERR     = "Unexpected HTTP Version"
	INVALID_VERB_ERR           = "Invalid HTTP Verb"
	INVALID_HEADER_NAME        = "Invalid Http field line key"
	INVALID_CONTENT_LENGTH     = "Invalid content length"
	CONTENT_LENGTH_MISMATCH    = "Content length mismatch with body"
)

const SEPERATOR = "\r\n"
const BUFFER_SIZE = 1024
const BUFFER_GROW_THRESHOLD = 32
const BUFFER_GROW_INCREMENT = 32

type parserState interface {
	parse(data []byte, parser *Parser) (int, error)
}

type Parser struct {
	Request    req.Request
	state      parserState
	end        bool
	processing bool
}

func (p *Parser) parse(data []byte) (int, error) {
	totConsumed := 0
	for consumed, err := p.state.parse(data, p); p.processing; consumed, err = p.state.parse(data[totConsumed:], p) {
		if err != nil {
			return totConsumed, err
		}
		totConsumed += consumed
	}
	return totConsumed, nil
}

type ErrorState struct{ message string }

func (e *ErrorState) parse(data []byte, parser *Parser) (int, error) {
	parser.processing = false
	parser.end = true
	return 0, fmt.Errorf("%s", e.message)
}

type DoneState struct{}

func (e *DoneState) parse(data []byte, parser *Parser) (int, error) {
	parser.processing = false
	parser.end = true
	return 0, nil
}

func (p *Parser) ParseFromReader(reader io.Reader) (*req.Request, error) {
	buffer := make([]byte, BUFFER_SIZE)
	bufferIndex := 0
	totalConsumed := 0

	p.Request = req.Request{}
	p.state = &ParseVerbState{}
	p.end = false

	for !p.end {
		// handle buffer resize if you are running out of space
		if len(buffer)-bufferIndex < BUFFER_GROW_THRESHOLD {
			buffer = append(buffer, make([]byte, BUFFER_GROW_INCREMENT)...)
		}

		// read into buffer
		size, err := reader.Read(buffer[bufferIndex:])
		if err != nil {
			// if err == io.EOF {
			// 	break
			// }
			return nil, err
		}
		bufferIndex += size

		// parse value
		p.processing = true
		consumed, err := p.parse(buffer[totalConsumed:bufferIndex])
		if err != nil {
			return nil, err
		}
		totalConsumed += consumed

	}

	return &p.Request, nil
}
