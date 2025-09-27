package parser

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeaders(t *testing.T) {

	// Test: Standard Headers
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}

	r, err := parser.ParseFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, "localhost:42069", r.Headers.Get("host"))
	assert.Equal(t, "curl/7.81.0", r.Headers.Get("user-agent"))
	assert.Equal(t, "*/*", r.Headers.Get("accept"))

	// Test: Empty Headers
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = parser.ParseFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r)
	assert.Equal(t, 0, len(r.Headers))

	// Test: Malformed Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	r, err = parser.ParseFromReader(reader)
	require.Error(t, err)

	// Test: missing end of Header
	reader = &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost localhost:42069\r\n",
		numBytesPerRead: 3,
	}
	r, err = parser.ParseFromReader(reader)
	require.Error(t, err)
}

func TestValidSingleHeader(t *testing.T) {
	// Test: Valid single header
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\nHost: localhost:42069\r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	r, err := parser.ParseFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r.Headers)
	assert.Equal(t, "localhost:42069", r.Headers.Get("Host"))
}

func TestValidSingleHeaderWithSpace(t *testing.T) {
	// Test: Valid single header with spacing
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n      Host: localhost:42069    \r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	r, err := parser.ParseFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r.Headers)
	assert.Equal(t, "localhost:42069", r.Headers.Get("Host"))
}

func TestValidMultipleHeaders(t *testing.T) {

	// Test: Valid multiple headers

	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n      Host: localhost:42069    \r\nTEST: testvalue\r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	r, err := parser.ParseFromReader(reader)

	require.NoError(t, err)
	require.NotNil(t, r.Headers)
	assert.Equal(t, "testvalue", r.Headers.Get("TEST"))
}

func TestCaseSensitiveNames(t *testing.T) {
	// Test: case sensitive header names
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n      host: localhost:42069    \r\nhOst: testvalue\r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	r, err := parser.ParseFromReader(reader)
	require.NoError(t, err)
	require.NotNil(t, r.Headers)
	assert.Equal(t, "localhost:42069, testvalue", r.Headers.Get("host"))
}

func TestInvalidHeaderSpacing(t *testing.T) {
	// Test: Invalid spacing header
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n       Host : localhost:42069       \r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	_, err := parser.ParseFromReader(reader)
	require.Error(t, err)
}

func TestInvalidHeaderName(t *testing.T) {
	// Test: Invalid spacing header
	reader := &chunkReader{
		data:            "GET / HTTP/1.1\r\n       Ho=t: localhost:42069       \r\n\r\n",
		numBytesPerRead: 3,
	}
	parser := &Parser{}
	_, err := parser.ParseFromReader(reader)
	require.Error(t, err)
}
