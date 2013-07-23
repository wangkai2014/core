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

type MiddlewareDummy2 struct {
	Middleware
}

func (mid *MiddlewareDummy2) Priority() int {
	return 5
}

func (mid *MiddlewareDummy2) Pre() {
	mid.C.Pub.Group.Set("result", "PRE2")
	mid.C.Cut()
}

func (mid *MiddlewareDummy2) Post() {
	mid.C.Pub.Group.Set("result", "POST2")
}

func TestMiddleware(t *testing.T) {
	c := &Core{
		App: NewApp(),
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

	mid = NewMiddlewares().Register(&MiddlewareDummy2{}, &MiddlewareDummy{}).Init(c)
	mid.Pre()

	if result() != "PRE2" {
		t.Fail()
	}

	mid.Post()

	if result() != "POST2" {
		t.Fail()
	}
}
