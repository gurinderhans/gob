package gob

import "reflect"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "io"
// import "fmt"

var (
  httpRequestType = reflect.TypeOf((*http.Request)(nil)).Elem()
  httpResponseWriterType = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
)

type Router struct {
  contextType reflect.Type
  routes []*route
  middlewares []reflect.Value
}

type route struct {
  Method string
  Path string
  Handler reflect.Value
}

func New(ctx interface{}) *Router {
  r := &Router{}
  r.contextType = reflect.TypeOf(ctx)
  return r
}

func (r *Router) Route(method, path string, fn interface{}) *Router {
  route := &route{
    Method: method,
    Path: path,
    Handler: reflect.ValueOf(fn),
  }
  r.routes = append(r.routes, route)

  return r
}

func (r *Router) Middleware(fn interface{}) *Router {
  r.middlewares = append(r.middlewares, reflect.ValueOf(fn))

  return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

  newCtx := reflect.New(r.contextType)

  /// execute all middleware first!
  for _, mw := range r.middlewares {
    res := mw.Call([]reflect.Value{
      newCtx,
      reflect.ValueOf(w),
      reflect.ValueOf(req),
    })

    if err := res[0].Interface(); err != nil {
      panic(err)
    }
  }

  for _, route := range r.routes {
    if route.Path == req.URL.Path && route.Method == req.Method {
      handlerType := route.Handler.Type()
      numParams := handlerType.NumIn()

      // FIXME: maybe move this logic to when adding routes ???
      if numParams == 2 {
        body, err := ioutil.ReadAll(req.Body)
        if err != nil {
          panic(err)
        }

        reqObj := reflect.New(handlerType.In(1).Elem()).Interface()
        if err := json.Unmarshal(body, &reqObj); err != nil {
          panic(err)
        }

        res := route.Handler.Call([]reflect.Value{
          newCtx,
          reflect.ValueOf(reqObj),
        })

        rsp, rspCode, rerr := res[0].Interface(), int(res[1].Int()), res[2].Interface()
        if rerr != nil {
          panic(rerr)
        }

        data, err := json.Marshal(rsp)
        w.WriteHeader(rspCode)
        io.WriteString(w, string(data))
      } else if numParams == 3 {
        if handlerType.In(1) == httpResponseWriterType && handlerType.In(2).Elem() == httpRequestType {
          route.Handler.Call([]reflect.Value{
            newCtx,
            reflect.ValueOf(w),
            reflect.ValueOf(req),
          })
        }
      }
    }
  }
}
