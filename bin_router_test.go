package core

import (
	"net/http"
	"testing"
)

func TestBinRouter(t *testing.T) {
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
			path:    "/",
			curpath: "",
			cut:     false,
		},
	}

	route := NewBinRouter().RootDir(pass)

	route.View(c)

	route.RootDir(fail).Register("blogpost", NewBinRouter().Register("example", pass))

	c.pri.path = "/blogpost/example"

	route.View(c)

	route.Register("blogpost", NewBinRouter().Register("example", fail))

	c.Pub.Errors.E404 = pass

	c.pri.path = "/blogpost/test"
	c.pri.curpath = ""

	route.View(c)
}
