package gob

import (
	"net/http"
)

type Request struct {
	*http.Request
	PathParams map[string]string
}
