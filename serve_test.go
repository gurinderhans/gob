package gob

import (
	"testing"
	// "net/http/httptest"
)

type BaseContext struct{}

func TestSimpleGet(t *testing.T) {
  router := NewRouter(BaseContext{}, "/api/v2")
  router.Get("/", fun
}
