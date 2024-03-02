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

	statusCode int
	bytes      int
}

func (b *BasicResponseWrapper) WriteHeader(code int) {
	if !b.wroteHeader {
		b.wroteHeader = true
		b.statusCode = code
		b.ResponseWriter.WriteHeader(code)
	}
}
