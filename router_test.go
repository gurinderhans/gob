package gob

import (
	"reflect"
	"testing"
)

var RootPrefix = "/api/v2"

type RootContext struct{}

func (c *RootContext) Foo(w ResponseWriter, req *Request) {}

type Req struct{}
type Rsp struct{}

type SubContext struct {
	*RootContext
}

func (c *SubContext) SubFoo(w ResponseWriter, req *Request) {}

func TestCreateNewRouter(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	if router.PathPrefix != RootPrefix {
		t.Errorf("Router has an invalid path prefix!")
		return
	}
	if router.root != router {
		t.Errorf("Incorrect root value set on router")
		return
	}
	if router.contextType != reflect.TypeOf(RootContext{}) {
		t.Errorf("Incorrect context type set for router")
		return
	}
}

func TestCreateSubrouter(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)

	subPrefix := "/sub"
	sub := router.Subrouter(SubContext{}, subPrefix)

	if sub.contextType != reflect.TypeOf(SubContext{}) {
		t.Errorf("Incorrect context set for subrouter")
		return
	}
	if sub.PathPrefix != RootPrefix+subPrefix {
		t.Errorf("Incorrect path prefix for sub router")
		return
	}
	if sub.parent != router {
		t.Errorf("Incorrect parent set for subrouter")
		return
	}
	if sub.root != router {
		t.Errorf("Incorrect root set for subrouter")
		return
	}
	if sub.routeTree != nil {
		t.Errorf("No subrouter should have a route tree initliazed")
		return
	}

	type SubSubContext struct {
		*SubContext
	}
	subsubRouter := sub.Subrouter(SubSubContext{}, "/down")
	if subsubRouter.contextType != reflect.TypeOf(SubSubContext{}) {
		t.Errorf("Incorrect context set for sub router")
		return
	}
	if subsubRouter.PathPrefix != RootPrefix+subPrefix+"/down" {
		t.Errorf("Incorrect path prefix set for sub router")
		return
	}
	if subsubRouter.parent != sub {
		t.Errorf("Incorrect parent set for subrouter")
		return
	}
	if subsubRouter.root != router {
		t.Errorf("Incorrect root set for sub router")
		return
	}
	if subsubRouter.routeTree != nil {
		t.Errorf("Subrouter should not have routeTree initialized!")
		return
	}
}

func TestInvalidRouterContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Creating router with wrong context did NOT panic!")
			return
		}
	}()

	type WrongContext []int
	NewRouter(WrongContext{}, RootPrefix)
}

func TestNilRouterContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Creating router with wrong context did NOT panic!")
			return
		}
	}()

	NewRouter(nil, RootPrefix)
}

func TestInvalidSubrouterContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Creating subrouter with missing parent context did NOT panic!")
			return
		}
	}()

	router := NewRouter(RootContext{}, RootPrefix)
	type MissingParentContext struct{}
	router.Subrouter(MissingParentContext{}, "/sub")
}

func TestNilSubrouterContext(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Creating subrouter with nil context did NOT panic!")
			return
		}
	}()

	router := NewRouter(RootContext{}, RootPrefix)
	router.Subrouter(nil, "/sub")
}

func TestAddHomeRoute(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Get("/", (*RootContext).Foo)

	if router.homeReceiver == nil {
		t.Errorf("Router's home receiver is NOT set!")
		return
	}
}

func TestAddSimpleRoute(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Get("/foo", (*RootContext).Foo)

	if router.homeReceiver != nil {
		t.Errorf("The home receiver on router should not bet set")
		return
	}
	res := router.routeTree.Find(RootPrefix + "/foo")
	if res == nil {
		t.Errorf("Added route not found")
		return
	}
	receiver := res.Value.(*requestReceiver)
	if receiver.Method != httpMethodGet {
		t.Errorf("Incorrect method set on route receiver")
		return
	}
	if receiver.Handler.Type() != reflect.TypeOf((*RootContext).Foo) {
		t.Errorf("Incorrect handler set on request receiver")
		return
	}
	if receiver.Router != router {
		t.Errorf("Incorrect router set on request receiver")
		return
	}
}

func TestAddComplexRoute(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Post("/user/:userId", (*RootContext).Foo)

	res := router.routeTree.Find(RootPrefix + "/user/me")
	if res == nil {
		t.Errorf("Added route not found!")
		return
	}

	userId, ok := res.Params["userId"]
	if !ok {
		t.Errorf("Param key does not exist!")
		return
	}
	if userId != "me" {
		t.Errorf("Incorrect value for param key")
		return
	}
}

func TestAddIncorrectContextRoute(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding route with incorrect context handler did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Get("/path", (*SubContext).SubFoo)
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding route handler with nil context did NOT panic!")
				return
			}
		}()
		router := NewRouter(RootContext{}, RootPrefix)
		router.Get("/path", func(w ResponseWriter, req *Request) {})
	}()
}

func TestAddComplexRouteHandler(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Get("/path", func(c *RootContext, req *Req) (*Rsp, int, error) {
		return &Rsp{}, 200, nil
	})
	res := router.root.routeTree.Find(RootPrefix + "/path")
	if res == nil {
		t.Errorf("Added route not found!")
		return
	}
}

func TestAddComplexInvalidRouteHandler(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Adding invalid complex route handler did NOT panic!")
			return
		}
	}()
	router := NewRouter(RootContext{}, RootPrefix)
	router.Get("/path", func(c *RootContext, req *Req) (int, error) {
		return 0, nil
	})
}

func TestInvalidRouteHandler(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Creating route with invalid handler did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Get("/path", []int{})
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Creating route with invalid handler did NOT panic!")
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Get("/path", func(c *RootContext, req *Req) (*Rsp, int, string) { return nil, 0, "" })
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Creating route with invalid handler did NOT panic!")
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Get("/path", func(c *RootContext, req *Req) {})
	}()
}

func TestRouterMiddleware(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Middleware(func(c *RootContext, w ResponseWriter, req *Request) error { return nil })
	if len(router.middlewares) != 1 {
		t.Errorf("Failed to add middleware to router")
		return
	}
}

func TestRouterInvalidMiddleware(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding invalid middleware did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Middleware([]int{})
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding middleware with incorrect context did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Middleware(func(c *SubContext, w ResponseWriter, req *Request) error { return nil })
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding middleware with no return error value did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Middleware(func(c *RootContext, w ResponseWriter, req *Request) {})
	}()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Adding middleware with incorrect return value did NOT panic!")
				return
			}
		}()

		router := NewRouter(RootContext{}, RootPrefix)
		router.Middleware(func(c *RootContext, w ResponseWriter, req *Request) string { return "" })
	}()
}
