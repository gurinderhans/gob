package gob

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"reflect"
)

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	requestPath := path.Clean(req.URL.Path)
	res := r.root.routeTree.Find(requestPath)
	if res == nil {
		fmt.Println("Path not found!")
		w.WriteHeader(404)
		return
	}

	var reqReceiver *requestReceiver
	switch res.Value.(type) {
	case *Router:
		router := res.Value.(*Router)
		reqReceiver = router.homeReceiver
	case *requestReceiver:
		reqReceiver = res.Value.(*requestReceiver)
	}

	if reqReceiver == nil || reqReceiver.Method != httpMethod(req.Method) {
		fmt.Println("404!")
		w.WriteHeader(404)
		return
	}

	request := &Request{req, res.Params}

	// create a router & context chain back upto root, run middlewares, then the final request handler!
	routerChain := []*Router{}
	for ro := reqReceiver.Router; ro != nil; ro = ro.parent {
		routerChain = append(routerChain, ro)
	}

	// FILL root context separately
	rootRouter := routerChain[len(routerChain)-1]
	parentContext := reflect.New(rootRouter.contextType)
	for _, mfn := range rootRouter.middlewares {
		res := mfn.Call([]reflect.Value{
			parentContext,
			reflect.ValueOf(w),
			reflect.ValueOf(request),
		})
		if err := res[0].Interface(); err != nil {
			fmt.Println("Error calling middleware!")
			w.WriteHeader(500)
			return
		}
	}

	// go down rest of the chain and fill the sub contexts!
	for i := len(routerChain) - 2; i >= 0; i-- {
		router := routerChain[i]
		ctx := reflect.New(router.contextType)

		field := reflect.Indirect(ctx).Field(0)
		field.Set(parentContext)

		// run middlewares for router
		for _, mfn := range router.middlewares {
			res := mfn.Call([]reflect.Value{
				parentContext,
				reflect.ValueOf(w),
				reflect.ValueOf(request),
			})
			if err := res[0].Interface(); err != nil {
				fmt.Println("Error calling middleware!")
				w.WriteHeader(500)
				return
			}
		}

		parentContext = ctx
	}

	handlerType := reqReceiver.Handler.Type()
	numParams := handlerType.NumIn()

	// FIXME: maybe move this logic to when adding routes ???
	if numParams <= 2 {
		var res []reflect.Value
		if numParams == 2 {
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			reqObj := reflect.New(handlerType.In(1).Elem()).Interface()
			if err := json.Unmarshal(body, &reqObj); err != nil {
				panic(err)
			}

			res = reqReceiver.Handler.Call([]reflect.Value{
				parentContext,
				reflect.ValueOf(reqObj),
			})
		} else {
			res = reqReceiver.Handler.Call([]reflect.Value{parentContext})
		}

		rsp, rspCode, rerr := res[0].Interface(), int(res[1].Int()), res[2].Interface()
		if rerr != nil {
			panic(rerr)
		}

		data, err := json.Marshal(rsp)
		if err != nil {
			panic(err)
		}
		w.WriteHeader(rspCode)
		io.WriteString(w, string(data))
	} else if numParams == 3 {
		if handlerType.In(1) == responseWriterType && handlerType.In(2).Elem() == requestType {
			reqReceiver.Handler.Call([]reflect.Value{
				parentContext,
				reflect.ValueOf(w),
				reflect.ValueOf(request),
			})
		}
	}
}
