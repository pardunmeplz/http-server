package response

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	req "tcpServer/internal/request"
)

type StatusCode string
type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(code StatusCode) error {

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
	return nil

}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Writer.Write(p)
}

const (
	OK                    StatusCode = "200"
	BAD_REQUEST           StatusCode = "400"
	INTERNAL_SERVER_ERROR StatusCode = "500"
)

func GetDefaultHeaders(contentLen int) req.Headers {
	output := make(req.Headers)

	output.Set("content-length", strconv.Itoa(contentLen))
	output.Set("connection", "close")
	output.Set("content-type", "text/plain")

	return output

}

func (w *Writer) WriteChunkedBody(resp *http.Response) {

	buffer := make([]byte, 1024)
	for {
		size, err := resp.Body.Read(buffer)

		if err == io.EOF {
			break
		} else if err != nil {
			log.Panicln("ERROR: ", err)
			return
		}

		_, err = w.WriteChunk(buffer[:size])
		if err != nil {
			log.Println("ERROR: ", err)
			return
		}
	}
	w.WriteChunkedBodyDone()
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
