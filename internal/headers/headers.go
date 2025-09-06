package headers

import (
	"fmt"
	"strings"
	req "tcpServer/internal/request"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (int, bool, error) {
	strVal := string(data)
	lineIndex := strings.Index(strVal, req.SEPERATOR)

	if lineIndex == -1 {
		return 0, false, nil
	}
	if lineIndex == 0 {
		return len(req.SEPERATOR), true, nil
	}

	nameEnd := strings.Index(strVal, ":")
	name := strings.TrimLeft(strVal[:nameEnd], " ")
	if strings.Index(name, " ") != -1 {
		return 0, false, fmt.Errorf("Malformed request headers %s", strVal)
	}
	h[name] = strings.Trim(strVal[nameEnd+1:lineIndex], " ")
	return lineIndex + len(req.SEPERATOR), false, nil
}

func NewHeaders() Headers {
	return make(Headers)
}
