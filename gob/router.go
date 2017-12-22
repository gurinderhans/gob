package gob

import (
	"path"
	"reflect"
)

type Router struct {
	PathPrefix string
	contextType  reflect.Type
	middlewares  []reflect.Value
	homeReceiver *requestReceiver
	root         *Router
	parent       *Router
	routeTree    *Trie
}

type requestReceiver struct {
	Handler reflect.Value
	Method  string
	Router  *Router
}

func NewRouter(ctx interface{}, basePath string) *Router {
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
	newRouter := &Router{
		PathPrefix:  path.Clean(r.PathPrefix + basePath),
		root:        r.root,
		parent:      r,
		contextType: reflect.TypeOf(ctx),
	}

	newRouter.root.routeTree.Add(newRouter.PathPrefix, newRouter)
	return newRouter
}

func (r *Router) Route(method, basePath string, fn interface{}) *Router {
	cleanedPath := path.Clean(r.PathPrefix + basePath)

	newReceiver := &requestReceiver{
		Router:  r,
		Method:  method,
		Handler: reflect.ValueOf(fn),
	}

	if cleanedPath == r.PathPrefix {
		r.homeReceiver = newReceiver
	} else {
		r.root.routeTree.Add(cleanedPath, newReceiver)
	}

	return r
}

func (r *Router) Middleware(fn interface{}) *Router {
	rfn := reflect.ValueOf(fn)
	if rfn.Type().In(0).Elem() == r.contextType {
		r.middlewares = append(r.middlewares, rfn)
	}
	return r
}
