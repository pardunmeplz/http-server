package request

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidSingleHeader(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 23, n)
	assert.False(t, done)
	data = data[n:]

	// Test: Done case
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 2, n)
	assert.True(t, done)
}

func TestValidSingleHeaderWithSpace(t *testing.T) {
	// Test: Valid single header with spacing
	headers := NewHeaders()
	data := []byte("      Host: localhost:42069    \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 33, n)
	assert.False(t, done)
}

func TestValidMultipleHeaders(t *testing.T) {

	// Test: Valid multiple headers
	headers := NewHeaders()
	data := []byte("      Host: localhost:42069    \r\nTEST: testvalue\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 33, n)
	assert.False(t, done)
	data = data[n:]

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "testvalue", headers.Get("TEST"))
	assert.Equal(t, 17, n)
	assert.False(t, done)
	data = data[n:]
}

func TestCaseSensitiveNames(t *testing.T) {
	// Test: case sensitive header names
	headers := NewHeaders()
	data := []byte("      host: localhost:42069    \r\nhOst: testvalue\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("Host"))
	assert.Equal(t, 33, n)
	assert.False(t, done)
	data = data[n:]

	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069, testvalue", headers.Get("host"))
	assert.Equal(t, 17, n)
	assert.False(t, done)
	data = data[n:]
}

func TestInvalidHeaderSpacing(t *testing.T) {
	// Test: Invalid spacing header
	headers := NewHeaders()
	data := []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestInvalidHeaderName(t *testing.T) {
	// Test: Invalid spacing header
	headers := NewHeaders()
	data := []byte("       Ho=t: localhost:42069       \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
