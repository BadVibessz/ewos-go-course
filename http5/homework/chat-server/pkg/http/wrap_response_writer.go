package http

import (
	"net/http"
)

type ResponseWriterWrapper interface {
	http.ResponseWriter

	Status() int
	BytesWritten() int
	Unwrap() http.ResponseWriter
}

type BasicResponseWrapper struct {
	http.ResponseWriter

	wroteHeader bool

	statusCode   int
	bytesWritten int
	response     string
}

func NewBasicResponseWrapper(rw http.ResponseWriter) *BasicResponseWrapper {
	return &BasicResponseWrapper{
		ResponseWriter: rw,
		wroteHeader:    false,
		statusCode:     0,
		bytesWritten:   0,
		response:       "",
	}
}

func (b *BasicResponseWrapper) Status() int {
	return b.statusCode
}

func (b *BasicResponseWrapper) BytesWritten() int {
	return b.bytesWritten
}

func (b *BasicResponseWrapper) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

func (b *BasicResponseWrapper) Response() string {
	return b.response
}

func (b *BasicResponseWrapper) WriteHeader(code int) {
	if !b.wroteHeader {
		b.wroteHeader = true
		b.statusCode = code
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *BasicResponseWrapper) Write(buf []byte) (int, error) {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}

	n, err := b.ResponseWriter.Write(buf)

	b.bytesWritten = n
	b.response = string(buf) // todo: if buf not representing a string?

	return n, err
}
