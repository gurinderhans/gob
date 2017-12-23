package gob

import (
	"net/http"
)

type ResponseWriter interface {
	http.ResponseWriter
}
