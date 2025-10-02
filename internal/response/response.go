package response

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	req "tcpServer/internal/request"
)

type StatusCode string
type ResponseState int

type Writer struct {
	Writer io.Writer
	State  ResponseState
}

func (w *Writer) WriteStatusLine(code StatusCode) error {
	if w.State != STATUS {
		return fmt.Errorf(INVALID_RESPONSE, "Status line")
	}

	w.State++
	switch code {
	case OK:
		_, err := w.Writer.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BAD_REQUEST:
		_, err := w.Writer.Write([]byte("HTTP/1.1 400 BAD REQUEST\r\n"))
		return err
	case INTERNAL_SERVER_ERROR:
		_, err := w.Writer.Write([]byte("HTTP/1.1 500 INTERNAL SERVER ERROR\r\n"))
		return err
	}
	return nil

}

func (w *Writer) WriteHeaders(headers req.Headers) error {
	if w.State != HEAD {
		return fmt.Errorf(INVALID_RESPONSE, "Headers")
	}
	for key := range headers {
		_, err := w.Writer.Write([]byte(key))
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte{':', ' '})
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte(headers.Get(key)))
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte{'\r', '\n'})
		if err != nil {
			return err
		}

	}
	w.Writer.Write([]byte{'\r', '\n'})
	w.State++
	return nil

}

func (w *Writer) WriteTrailers(headers req.Headers) error {
	if w.State != TRAILER {
		return fmt.Errorf(INVALID_RESPONSE, "Trailers")
	}

	for key := range headers {
		_, err := w.Writer.Write([]byte(key))
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte{':', ' '})
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte(headers.Get(key)))
		if err != nil {
			return err
		}
		_, err = w.Writer.Write([]byte{'\r', '\n'})
		if err != nil {
			return err
		}

	}
	w.Writer.Write([]byte{'\r', '\n'})
	w.State++
	return nil

}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.State != BODY {
		return 0, fmt.Errorf(INVALID_RESPONSE, "Body")
	}
	w.State++
	return w.Writer.Write(p)
}

const (
	OK                    StatusCode = "200"
	BAD_REQUEST           StatusCode = "400"
	INTERNAL_SERVER_ERROR StatusCode = "500"
)

const (
	STATUS  ResponseState = 1
	HEAD    ResponseState = 2
	BODY    ResponseState = 3
	TRAILER ResponseState = 4
	DONE    ResponseState = 5
)

const (
	INVALID_RESPONSE = "Can not write %s in response"
)

func GetDefaultHeaders(contentLen int) req.Headers {
	output := make(req.Headers)

	output.Set("content-length", strconv.Itoa(contentLen))
	output.Set("connection", "close")
	output.Set("content-type", "text/plain")

	return output

}

func (w *Writer) WriteChunkedBody(resp *http.Response) (int, error) {
	if w.State != BODY {
		return 0, fmt.Errorf(INVALID_RESPONSE, "Body")
	}

	buffer := make([]byte, 1024)
	consumed := 0
	for {
		size, err := resp.Body.Read(buffer)

		if err == io.EOF {
			break
		} else if err != nil {
			return consumed, err
		}

		_, err = w.WriteChunk(buffer[:size])
		consumed += size
		if err != nil {
			return consumed, err
		}
	}
	w.WriteChunkedBodyDone()
	w.State++
	return consumed, nil
}

func (w *Writer) WriteChunk(p []byte) (int, error) {
	// write hex num /r/n
	sizeA, err := w.Writer.Write([]byte(fmt.Sprintf("%x\r\n", len(p))))
	sizeB, err := w.Writer.Write(p)
	sizeC, err := w.Writer.Write([]byte{'\r', '\n'})
	return sizeA + sizeB + sizeC, err

}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte{'0', '\r', '\n', '\r', '\n'})
}
