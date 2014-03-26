package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MiddlewareDummy struct {
	Middleware
}

func (mid *MiddlewareDummy) Pre() {
	mid.C.Pub.Group.Set("result", "PRE")
}

func (mid *MiddlewareDummy) Post() {
	mid.C.Pub.Group.Set("result", "POST")
}

type MiddlewareDummy2 struct {
	Middleware
}

func (mid *MiddlewareDummy2) Priority() int {
	return 5
}

func (mid *MiddlewareDummy2) Pre() {
	mid.C.Pub.Group.Set("result", "PRE2")
	mid.C.Terminate()
}

func (mid *MiddlewareDummy2) Post() {
	mid.C.Pub.Group.Set("result", "POST2")
}

func TestMiddleware(t *testing.T) {
	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Context) {
		result := func() string {
			return c.Pub.Group.Get("result")
		}

		mid := NewMiddlewares().Register(&MiddlewareDummy{}).Init(c)
		mid.Pre()

		if result() != "PRE" {
			t.Fail()
		}

		mid.Post()

		if result() != "POST" {
			t.Fail()
		}

		mid = NewMiddlewares().Register(&MiddlewareDummy2{}, &MiddlewareDummy{}).Init(c)
		mid.Pre()

		if result() != "PRE2" {
			t.Fail()
		}

		mid.Post()

		if result() != "POST2" {
			t.Fail()
		}
	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
