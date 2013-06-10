package core

import (
	"fmt"
	"net/http"
	"testing"
)

func TestVhost(t *testing.T) {
	pass := RouteHandlerFunc(func(c *Core) {
		// Do nothing, because it's a pass
	})

	fail := RouteHandlerFunc(func(c *Core) {
		t.Fail()
	})

	c := &Core{
		Pub: Public{
			Errors: Errors{
				E403: fail,
				E404: fail,
				E500: fail,
			},
		},
		Req: &http.Request{
			Host:   "example.com:1234",
			Header: http.Header{},
		},
		pri: private{
			cut: false,
		},
	}

	hosts := NewVHost(Map{
		`example.com`: pass,
	})

	hosts.View(c)

	c.Req.Host = "www.example.com:1234"

	c.Pub.Errors.E404 = pass

	hosts = NewVHost(Map{
		`example.com`: fail,
	})

	hosts.View(c)
}

func TestVhostRegExp(t *testing.T) {
	possible_pass := RouteHandlerFunc(func(c *Core) {
		if c.Pub.Group.Get("subdomain") != "hello" {
			t.Fail()
			fmt.Println(c.Req.Host)
		}
	})

	pass := RouteHandlerFunc(func(c *Core) {
		// Do nothing, because it's a pass
	})

	fail := RouteHandlerFunc(func(c *Core) {
		t.Fail()
		fmt.Println(c.Req.Host)
	})

	c := &Core{
		Pub: Public{
			Errors: Errors{
				E403: fail,
				E404: fail,
				E500: fail,
			},
			Group: Group{},
		},
		Req: &http.Request{
			Host:   "hello.example.com:1234",
			Header: http.Header{},
		},
		pri: private{
			cut: false,
		},
	}

	const rule = `^(?P<subdomain>[a-z]+)\.example\.com`

	vhost := NewVHostRegExp(Map{
		rule: possible_pass,
	})

	vhost.View(c)

	c.Req.Host = "www.example.com:1234"
	c.Pub.Group = Group{}

	vhost = NewVHostRegExp(Map{
		rule: pass,
	})

	vhost.View(c)

	c.Req.Host = "www.hello.com:1234"

	c.Pub.Errors.E404 = pass

	vhost = NewVHostRegExp(Map{
		rule: fail,
	})

	vhost.View(c)
}
