package core

import (
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
)

func TestSession(t *testing.T) {
	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Session().Set("hello world")
	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	client := &http.Client{}

	client.Jar, _ = cookiejar.New(nil)

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		if c.Session().Get().(string) != "hello world" {
			t.Fail()
		}
	})

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Session().Destroy()
	})

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		if c.Session().Get() != nil {
			t.Fail()
		}
	})

	client.Get(ts.URL)
}

func TestSessionAdv(t *testing.T) {
	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Session().Adv().Set("world", "hello!")
		c.Session().Adv().Save()
	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	client := &http.Client{}

	client.Jar, _ = cookiejar.New(nil)

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		if c.Session().Adv().Get("world").(string) != "hello!" {
			t.Fail()
		}
	})

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		c.Session().Adv().Set("world", nil)
		c.Session().Adv().Save()
	})

	client.Get(ts.URL)

	App.TestView = RouteHandlerFunc(func(c *Context) {
		if c.Session().Adv().Get("world") != nil {
			t.Fail()
		}
	})

	client.Get(ts.URL)
}
