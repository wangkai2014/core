package core

import (
	"net/http"
	"testing"
)

func TestRouter(t *testing.T) {
	possible_pass := RouteHandlerFunc(func(c *Core) {
		if c.Pub.Group.Get("title") != "test" {
			t.Fail()
		}

		if c.Pub.Group.GetInt("id") != 5 {
			t.Fail()
		}
	})

	pass := RouteHandlerFunc(func(c *Core) {
		// Do nothing, it's an automactic pass!
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
			Group: Group{},
		},
		Req: &http.Request{
			Header: http.Header{},
		},
		pri: private{
			path:    "/blogpost/test-5",
			curpath: "",
			cut:     false,
		},
	}

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
}
