package core

import (
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

func TestMiddleware(t *testing.T) {
	c := &Core{
		Pub: Public{
			Group: Group{},
		},
		pri: private{
			cut: false,
		},
	}

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
}
