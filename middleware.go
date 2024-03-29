package core

import (
	"reflect"
	"sort"
	"sync"
)

// Middleware Interface
type MiddlewareInterface interface {
	Init(*Context)
	Pre()
	Post()
	Priority() int
	getType() reflect.Type
	setType(reflect.Type)
}

// Implement MiddlewareInterface
type Middleware struct {
	C  *Context
	_t reflect.Type
	_s sync.RWMutex
}

// Init
func (mid *Middleware) Init(c *Context) {
	mid.C = c
}

// Pre boot
func (mid *Middleware) Pre() {
	// Do nothing
}

// Post boot
func (mid *Middleware) Post() {
	// Do nothing
}

// Priority
func (mid *Middleware) Priority() int {
	return 10
}

func (mid *Middleware) getType() reflect.Type {
	mid._s.RLock()
	defer mid._s.RUnlock()
	return mid._t
}

func (mid *Middleware) setType(t reflect.Type) {
	mid._s.Lock()
	defer mid._s.Unlock()
	mid._t = t
}

type _middlewares []MiddlewareInterface

func (mid _middlewares) Len() int {
	return len(mid)
}

func (mid _middlewares) Less(i, j int) bool {
	return mid[i].Priority() < mid[j].Priority()
}

func (mid _middlewares) Swap(i, j int) {
	mid[i], mid[j] = mid[j], mid[i]
}

type Middlewares struct {
	sync.Mutex
	items _middlewares
	c     *Context
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
		mid.items = _middlewares{}
	}
	mid.items = append(mid.items, middlewares...)
	sort.Sort(mid.items)
	return mid
}

// Clear Middlewares
func (mid *Middlewares) Clear() *Middlewares {
	if mid.c != nil {
		return mid
	}
	mid.items = nil
	return mid
}

// Init Middlewares, return initialised structure.
func (mid *Middlewares) Init(c *Context) *Middlewares {
	if mid.c != nil || !c.App.MiddlewareEnabled {
		return mid
	}
	middlewares := NewMiddlewares()

	middlewares.items = _middlewares{}
	middlewares.c = c
	for _, middleware := range mid.items {
		t := middleware.getType()
		if t == nil {
			t = reflect.Indirect(reflect.ValueOf(middleware)).Type()
			middleware.setType(t)
		}

		newmiddleware := reflect.New(t).Interface().(MiddlewareInterface)
		newmiddleware.Init(c)
		middlewares.items = append(middlewares.items, newmiddleware)
	}
	return middlewares
}

// Pre boot
func (mid *Middlewares) Pre() {
	if mid.c == nil {
		return
	}
	for key, middleware := range mid.items {
		middleware.Pre()
		if mid.c.Terminated() {
			mid.items = mid.items[:key+1]
			return
		}
	}
}

// Post boot, you may want to use the keyword 'defer'
// Execute in Reverse unlike Pre().
func (mid *Middlewares) Post() {
	if mid.c == nil {
		return
	}
	for i := len(mid.items); i > 0; i-- {
		mid.items[i-1].Post()
	}
}
