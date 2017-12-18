package main

// import "fmt"
import "net/http"
import "io"
import "errors"
import "./gob"

type RequestContext struct {
  Writer http.ResponseWriter
}

func (c *RequestContext) RawHandler(w http.ResponseWriter, req *http.Request) {
  // here you should be able to read the request object byte by byte....
  // and write to the response object byte by byte....
  w.Header().Set("random", "stuff")
  w.WriteHeader(203)
	io.WriteString(w, "Raw handler, Hello world!")
}

func (c *RequestContext) SetContext(w http.ResponseWriter, req *http.Request) error {
  c.Writer=w
  return nil
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
    Users: []string{"user one", "user 2"},
  }

  if cReq.Name == "rnad" {
    err := errors.New("rnad name is not allowed!")
    return nil, 400, err
  }

  c.Writer.Header().Set("Location", "https://www.api.com/chat/z87dadsd8")

  return rsp, 201, nil
}

func main() {
  router := gob.New(RequestContext{}).
    Middleware((*RequestContext).SetContext).
    Route("GET", "/raw", (*RequestContext).RawHandler).
    Route("POST","/chat", (*RequestContext).CreateChat)

  http.ListenAndServe(":3000", router)
}
