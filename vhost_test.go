package core

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVhost(t *testing.T) {
	App := NewApp()

	App.Debug = true

	pass := RouteHandlerFunc(func(c *Context) {
		// Do nothing, because it's a pass
	})

	fail := RouteHandlerFunc(func(c *Context) {
		t.Fail()
	})

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Pub.Errors.E403 = fail
		c.Pub.Errors.E404 = fail
		c.Pub.Errors.E500 = fail

		c.Req.Host = "example.com:1234"

		hosts := NewVHost().Register(Map{
			`example.com`: pass,
		})

		hosts.View(c)

		c.Req.Host = "www.example.com:1234"

		c.Pub.Errors.E404 = pass

		hosts = NewVHost().Register(Map{
			`example.com`: fail,
		})

		hosts.View(c)

	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}

func TestVhostRegExp(t *testing.T) {
	App := NewApp()

	App.Debug = true

	possible_pass := RouteHandlerFunc(func(c *Context) {
		if c.Pub.Group.Get("subdomain") != "hello" {
			t.Fail()
			fmt.Println(c.Req.Host)
		}
	})

	pass := RouteHandlerFunc(func(c *Context) {
		// Do nothing, because it's a pass
	})

	fail := RouteHandlerFunc(func(c *Context) {
		t.Fail()
		fmt.Println(c.Req.Host)
	})

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Pub.Errors.E403 = fail
		c.Pub.Errors.E404 = fail
		c.Pub.Errors.E500 = fail

		c.Req.Host = "hello.example.com:1234"

		const rule = `^(?P<subdomain>[a-z]+)\.example\.com`

		vhost := NewVHostRegExp().Register(Map{
			rule: possible_pass,
		})

		vhost.View(c)

		c.Req.Host = "www.example.com:1234"
		c.Pub.Group = Group{}

		vhost = NewVHostRegExp().Register(Map{
			rule: pass,
		})

		vhost.View(c)

		c.Req.Host = "www.hello.com:1234"

		c.Pub.Errors.E404 = pass

		vhost = NewVHostRegExp().Register(Map{
			rule: fail,
		})

		vhost.View(c)

	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
