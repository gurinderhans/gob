package main

import "fmt"
import "net/http"
import "io"
import "errors"
import "./gob"

type RequestContext struct {
  RandomVal int
}

func (c *RequestContext) RawHandler(w http.ResponseWriter, req *http.Request) {
  // here you should be able to read the request object byte by byte....
  // and write to the response object byte by byte....
  w.WriteHeader(203)
	io.WriteString(w, "Raw handler, Hello world!")
}

func (c *RequestContext) SetContext(w http.ResponseWriter, req *http.Request) {
  c.RandomVal = 3
  fmt.Println("set context")
}

type ChatReq struct {
  Name string `json:"name"`
  UserIDs []string `json:"users"`
}

type ChatResp struct {
  Name string `json:"name"`
  Users []string `json:"users"`
}

func (c *RequestContext) CreateChat(cReq *ChatReq) (*ChatResp, int, error) {
  rsp := &ChatResp{
    Name: "some chat",
    Users: []string{"new user "+string(c.RandomVal), cReq.UserIDs[0]+"_ name"},
  }

  if cReq.Name == "rnad" {
    return nil, 400, errors.New("rnad name not allowed!")
  }

  // ??? c.Headers["Location"] = "https://www.google.com"

  return rsp, 301, nil
}

func main() {
  router := gob.New(RequestContext{})

  router.Middleware((*RequestContext).SetContext)

  router.Route("GET", "/raw", (*RequestContext).RawHandler)
  router.Route("POST","/chat", (*RequestContext).CreateChat)

  http.ListenAndServe(":3000", router)
}
