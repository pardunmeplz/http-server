package response

import (
	"io"
	"strconv"
	"tcpServer/internal/headers"
)

type StatusCode string

const (
	OK                    StatusCode = "200"
	BAD_REQUEST           StatusCode = "400"
	INTERNAL_SERVER_ERROR StatusCode = "500"
)

func WriteStatusLine(writer io.Writer, code StatusCode) error {

	switch code {
	case OK:
		_, err := writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BAD_REQUEST:
		_, err := writer.Write([]byte("HTTP/1.1 400 BAD REQUEST\r\n"))
		return err
	case INTERNAL_SERVER_ERROR:
		_, err := writer.Write([]byte("HTTP/1.1 500 INTERNAL SERVER ERROR\r\n"))
		return err
	}
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	output := make(headers.Headers)

	output.Set("content-length", strconv.Itoa(contentLen))
	output.Set("connection", "close")
	output.Set("content-type", "text/plain")

	return output

}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	for key := range headers {
		_, err := w.Write([]byte(key))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{':', ' '})
		if err != nil {
			return err
		}
		_, err = w.Write([]byte(headers.Get(key)))
		if err != nil {
			return err
		}
		_, err = w.Write([]byte{'\r', '\n'})
		if err != nil {
			return err
		}

	}
	w.Write([]byte{'\r', '\n'})
	return nil
}
