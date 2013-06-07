package core

import (
	"reflect"
	"sync"
)

// Middleware Interface
type MiddlewareInterface interface {
	Init(*Core)
	Html()
	Pre()
	Post()
}

// Implement MiddlewareInterface
type Middleware struct {
	C *Core
}

// Init
func (mid *Middleware) Init(c *Core) {
	mid.C = c
}

// Html
func (mid *Middleware) Html() {

}

// Pre boot
func (mid *Middleware) Pre() {
	// Do nothing
}

// Post boot
func (mid *Middleware) Post() {
	// Do nothing
}

type Middlewares struct {
	sync.Mutex
	items []MiddlewareInterface
	c     *Core
}

// Construct New Middleware
func NewMiddlewares() *Middlewares {
	return &Middlewares{}
}

// Register Middlewares
func (mid *Middlewares) Register(middlewares ...MiddlewareInterface) *Middlewares {
	if mid.c == nil {
		mid.Lock()
		defer mid.Unlock()
	}
	if mid.items == nil {
		mid.items = []MiddlewareInterface{}
	}
	mid.items = append(mid.items, middlewares...)
	return mid
}

func (mid *Middlewares) getItems() []MiddlewareInterface {
	mid.Lock()
	defer mid.Unlock()
	return append([]MiddlewareInterface{}, mid.items...)
}

// Init Middlewares, return initialised structure.
func (mid *Middlewares) Init(c *Core) *Middlewares {
	if mid.c != nil {
		return mid
	}
	middlewares := NewMiddlewares()
	middlewares.c = c
	for _, middleware := range mid.getItems() {
		newmiddleware := reflect.New(
			reflect.Indirect(reflect.ValueOf(middleware)).Type()).Interface().(MiddlewareInterface)
		newmiddleware.Init(c)
		middlewares.Register(newmiddleware)
	}
	return middlewares
}

// Html
func (mid *Middlewares) Html() {
	if mid.c == nil {
		return
	}
	for _, middleware := range mid.items {
		middleware.Html()
		if mid.c.CutOut() {
			return
		}
	}
}

// Pre boot
func (mid *Middlewares) Pre() {
	if mid.c == nil {
		return
	}
	for _, middleware := range mid.items {
		middleware.Pre()
		if mid.c.CutOut() {
			return
		}
	}
}

// Post boot, you may want to use the keyword 'defer'
func (mid *Middlewares) Post() {
	if mid.c == nil {
		return
	}
	for _, middleware := range mid.items {
		middleware.Post()
	}
}
