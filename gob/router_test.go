package gob

import (
	"net/http"
	"path"
	"reflect"
	"testing"
  // "fmt"
)

type RootContext struct{}

func (c *RootContext) Foo(w http.ResponseWriter, req *Request) {
}

type SubContext struct{}

var RootPrefix = "/api/v2"

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
	if sub.PathPrefix != path.Clean(RootPrefix+subPrefix) {
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

	type SubSubContext struct{}
	subsubRouter := sub.Subrouter(SubSubContext{}, "/down")
	if subsubRouter.contextType != reflect.TypeOf(SubSubContext{}) {
		t.Errorf("Incorrect context set for sub router")
		return
	}
	if subsubRouter.PathPrefix != path.Clean(RootPrefix+subPrefix+"/down") {
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

func TestAddHomeRoute(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Route("GET", "/", (*RootContext).Foo)

	if router.homeReceiver == nil {
		t.Errorf("Router's home receiver is NOT set!")
		return
	}
	if router.homeReceiver.Method != "GET" {
		t.Errorf("Incorrect method set on receiver!")
		return
	}
	if router.homeReceiver.Router != router {
		t.Errorf("Incorrect router set on home receiver")
		return
	}
	if router.homeReceiver.Handler.Type() != reflect.TypeOf((*RootContext).Foo) {
		t.Errorf("Incorrect handler set on home receiver")
		return
	}
}

func TestAddRoutes(t *testing.T) {
	router := NewRouter(RootContext{}, RootPrefix)
	router.Route("GET", "/foo", (*RootContext).Foo)

  res := router.routeTree.Find(RootPrefix + "/foo")
  if res == nil {
    t.Errorf("Route was added incorrectly!")
    return
  }
  receiver := res.Value.(*requestReceiver)
  if receiver.Method != "GET" {
    t.Errorf("Incorrect method set on route")
    return
  }
	if receiver.Handler.Type() != reflect.TypeOf((*RootContext).Foo) {
		t.Errorf("Incorrect handler set on home receiver")
		return
	}
  if receiver.Router != router {
    t.Errorf("Incorrect router set on receiver!")
    return
  }

  router.Route("POST", "/user/:userId", (*RootContext).Foo)

  res = router.routeTree.Find(RootPrefix + "/user/me")
  if res == nil {
    t.Errorf("Unable to match route!")
    return
  }
  pathParams := res.Params
  id, ok := pathParams["userId"]
  if !ok {
    t.Errorf("Key does not exist!")
    return
  }
  if id != "me" {
    t.Errorf("Wrong value for key!")
    return
  }
}
