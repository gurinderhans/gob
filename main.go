package main

import (
	"./gob"
	// "errors"
	// "io"
	"net/http"
	// "fmt"
)

type C struct {}

type Req struct {}

func (c *C) Foo(r *Req) (struct{}, int, error) {
  return struct{}{}, 200, nil
}

// -> both are same paths
// /p/:id
// /p/:man
//
// /p/:id/comments
// /p/:user/friends
// -> converts to
//   user
// /p/:/comments
//   user
// /p/:/friends
func main() {
  router := gob.NewRouter(C{}, "/")
  router.Post("/", (*C).Foo)
	http.ListenAndServe(":3000", router)
}
