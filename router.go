package gob

import (
	"path"
	"reflect"
)

type httpMethod string

const (
	httpMethodGet    = httpMethod("GET")
	httpMethodPost   = httpMethod("POST")
	httpMethodPut    = httpMethod("PUT")
	httpMethodDelete = httpMethod("DELETE")
)

var (
	responseWriterType = reflect.TypeOf((*ResponseWriter)(nil)).Elem()
	requestType        = reflect.TypeOf((*Request)(nil)).Elem()
	errorType          = reflect.TypeOf((*error)(nil)).Elem()
)

type Router struct {
	PathPrefix   string
	contextType  reflect.Type
	middlewares  []reflect.Value
	homeReceiver *requestReceiver
	root         *Router
	parent       *Router
	routeTree    *Trie
}

type requestReceiver struct {
	Handler reflect.Value
	Method  httpMethod
	Router  *Router
}

func NewRouter(ctx interface{}, basePath string) *Router {
	validateContext(ctx, nil)

	newRouter := &Router{
		PathPrefix:  path.Clean(basePath),
		contextType: reflect.TypeOf(ctx),
		routeTree:   NewTrie(),
	}

	newRouter.root = newRouter
	newRouter.root.routeTree.Add(newRouter.PathPrefix, newRouter)
	return newRouter
}

func (r *Router) Subrouter(ctx interface{}, basePath string) *Router {
	validateContext(ctx, r.contextType)

	newRouter := &Router{
		PathPrefix:  path.Clean(r.PathPrefix + basePath),
		root:        r.root,
		parent:      r,
		contextType: reflect.TypeOf(ctx),
	}

	newRouter.root.routeTree.Add(newRouter.PathPrefix, newRouter)
	return newRouter
}

func (r *Router) Middleware(fn interface{}) *Router {
	vfn := reflect.ValueOf(fn)

	validateMiddlewareHandler(vfn, r.contextType)

	r.middlewares = append(r.middlewares, vfn)
	return r
}

func (r *Router) Get(path string, fn interface{}) *Router {
	return r.newReceiver(httpMethodGet, path, fn)
}

func (r *Router) Post(path string, fn interface{}) *Router {
	return r.newReceiver(httpMethodPost, path, fn)
}

func (r *Router) Put(path string, fn interface{}) *Router {
	return r.newReceiver(httpMethodPut, path, fn)
}

func (r *Router) Delete(path string, fn interface{}) *Router {
	return r.newReceiver(httpMethodDelete, path, fn)
}

func (r *Router) newReceiver(method httpMethod, basePath string, fn interface{}) *Router {
	vfn := reflect.ValueOf(fn)

	validateHandler(vfn, r.contextType)

	newReceiver := &requestReceiver{
		Router:  r,
		Method:  method,
		Handler: vfn,
	}
	fullPath := path.Clean(r.PathPrefix + basePath)
	if r.PathPrefix == fullPath {
		r.homeReceiver = newReceiver
	} else {
		r.root.routeTree.Add(fullPath, newReceiver)
	}
	return r
}

func validateContext(ctx interface{}, parentCtxType reflect.Type) {
	ctxType := reflect.TypeOf(ctx)
	if ctxType.Kind() != reflect.Struct {
		panic("Context needs to be a struct type")
	}

	if parentCtxType != nil && parentCtxType != ctxType {
		if ctxType.NumField() == 0 {
			panic("Context needs to have first field be a pointer to parent context")
		}
		fldType := ctxType.Field(0).Type
		if fldType != reflect.PtrTo(parentCtxType) {
			panic("Context needs to have first field be a pointer to parent context")
		}
	}
}

func validateHandler(vfn reflect.Value, ctxType reflect.Type) {
	vType := vfn.Type()
	if vType.Kind() != reflect.Func {
		panic("Handler type should be a func.")
	}

	numIn := vType.NumIn()
	if numIn == 0 {
		panic("Handler's first param should be a context param")
	} else if vType.In(0).Elem() != ctxType {
		panic("Invalid handler, first param context doesn't match router context")
	}

	numOut := vType.NumOut()
	if numOut != 0 && numOut != 3 {
		panic("Invalid number of out parameters")
	} else if numOut == 0 {
		if vType.In(1) != responseWriterType || vType.In(2).Elem() != requestType {
			panic("Invalid handler, should contain response writer and request type as params")
		}
	} else if numOut == 3 {
		if vType.Out(1).Kind() != reflect.Int || vType.Out(2) != errorType {
			panic("Invalid handler, should have int and error as out param types")
		}
	}
}

func validateMiddlewareHandler(mfn reflect.Value, ctxType reflect.Type) {
	mType := mfn.Type()
	if mType.Kind() != reflect.Func {
		panic("Handler type should be a func.")
	}

	numIn := mType.NumIn()
	if numIn == 0 {
		panic("Middleware's first param should be a context param")
	} else if mType.In(0).Elem() != ctxType {
		panic("Invalid handler, context param doesn't match router context")
	}

	numOut := mType.NumOut()
	if numOut == 0 {
		panic("Middlewares need to return an error type")
	} else if mType.Out(0) != errorType {
		panic("Return type for middleware isn't an error type")
	}
}
