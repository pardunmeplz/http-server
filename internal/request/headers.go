package request

import (
	"fmt"
	"strings"
	"unicode"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (int, bool, error) {
	strVal := string(data)
	lineIndex := strings.Index(strVal, SEPERATOR)

	// handle incomplete line and done cases
	if lineIndex == -1 {
		return 0, false, nil
	}
	if lineIndex == 0 {
		return len(SEPERATOR), true, nil
	}

	// handle header name
	nameEnd := strings.Index(strVal, ":")
	if nameEnd == -1 {
		return 0, false, nil
	}
	name := strings.TrimLeft(strVal[:nameEnd], " ")
	if !isValidFieldName(name) {
		return 0, false, fmt.Errorf("Malformed request headers %s", strVal)
	}
	name = strings.ToLower(name)

	// handle header value
	val := strings.Trim(strVal[nameEnd+1:lineIndex], " ")
	if _, ok := h[name]; ok {
		val = h[name] + ", " + val
	}
	h[strings.ToLower(name)] = val

	return lineIndex + len(SEPERATOR), false, nil
}

func (h Headers) Get(name string) string {
	return h[strings.ToLower(name)]
}

func (h Headers) Set(name string, value string) {
	h[strings.ToLower(name)] = value
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

func NewHeaders() Headers {
	return make(Headers)
}
