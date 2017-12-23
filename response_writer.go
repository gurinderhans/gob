package gob

import (
	"net/http"
)

type ResponseWriter struct {
	http.ResponseWriter
}

func (w ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w ResponseWriter) Write(data []byte) (int, error) {
	return w.ResponseWriter.Write(data)
}

func (w ResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}
