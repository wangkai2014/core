package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	App := NewApp()

	App.Debug = true

	possible_pass := RouteHandlerFunc(func(c *Context) {
		if c.Pub.Group.Get("title") != "test" {
			t.Fail()
		}

		if c.Pub.Group.GetInt("id") != 5 {
			t.Fail()
		}
	})

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

		c.pri.path = "/blogpost/test-5"

		route := NewRouter().RegisterMap(Map{
			`^/blogpost`: NewRouter().RegisterMap(Map{
				`^/(?P<title>[a-z]+)-(?P<id>\d+)/?$`: possible_pass,
			}),
		})

		route.Load(c)

		c.pri.path = "/55"
		c.pri.curpath = ""

		c.Pub.Errors.E404 = pass

		route.Load(c)

	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
