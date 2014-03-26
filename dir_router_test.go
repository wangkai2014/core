package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDirRouter(t *testing.T) {
	App := NewApp()

	App.Debug = true

	pass := RouteHandlerFunc(func(c *Context) {
		// Do nothing, it's an automactic pass!
	})

	fail := RouteHandlerFunc(func(c *Context) {
		t.Fail()
	})

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Pub.Errors.E403 = fail
		c.Pub.Errors.E404 = fail
		c.Pub.Errors.E500 = fail

		route := NewDirRouter().RootDir(pass)

		route.View(c)

		route.RootDir(fail).Register("blogpost", NewDirRouter().Register("example", pass))

		c.pri.path = "/blogpost/example"

		route.View(c)

		route.Register("blogpost", NewDirRouter().Register("example", fail))

		c.Pub.Errors.E404 = pass

		c.pri.path = "/blogpost/test"
		c.pri.curpath = ""

		route.View(c)

		c.pri.path = "/blogpost/example/test"
		c.pri.curpath = ""

		route.View(c)
	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
