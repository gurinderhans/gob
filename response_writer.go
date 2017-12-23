package gob

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
}

type responseWriter struct {
	http.ResponseWriter
}

func (w responseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w responseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
