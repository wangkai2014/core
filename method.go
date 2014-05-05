package core

import (
	"reflect"
	"strings"
	"sync"
)

func execMethodInterface(c *Context, me MethodInterface) {
	t := me.getType()
	if t == nil {
		t = reflect.Indirect(reflect.ValueOf(me)).Type()
		me.setType(t)
	}

	vc := reflect.New(t)

	view := vc.MethodByName("View")
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(c)
	view.Call(in)

	c.Auto().PopulateStructFieldsValue(vc, "C")

	if c.Terminated() {
		return
	}

	in = make([]reflect.Value, 0)
	method := vc.MethodByName("Prepare")
	method.Call(in)

	if c.Terminated() {
		return
	}

	is := c.Is()

	if is.WebSocketRequest() {
		method = vc.MethodByName("Ws")
		method.Call(in)
		if c.Terminated() {
			goto finish
		}
	}

	switch c.Req.Method {
	case "GET", "HEAD", "POST":
		// Do nothing
	default:
		goto requestDealer
	}

	if is.AjaxRequest() {
		method = vc.MethodByName("Ajax")
		method.Call(in)
		if c.Terminated() {
			goto finish
		}
	}

requestDealer:

	switch c.Req.Method {
	case "GET", "HEAD":
		method = vc.MethodByName("Get")
		method.Call(in)
	case "POST", "DELETE", "PUT", "PATCH", "OPTIONS":
		method = vc.MethodByName(strings.Title(strings.ToLower(c.Req.Method)))
		method.Call(in)
	}

finish:

	method = vc.MethodByName("Finish")
	method.Call(in)
}

type MethodInterface interface {
	View(*Context)
	Prepare()
	Ws()
	Ajax()
	Get()
	Post()
	Delete()
	Put()
	Patch()
	Options()
	Finish()
	getType() reflect.Type
	setType(reflect.Type)

	asn_Core_0001() // Assert Serial Number
}

type Method struct {
	C  *Context `json:"-" xml:"-"`
	_t reflect.Type
	_s sync.RWMutex
}

func (me *Method) View(c *Context) {
	me.C = c
}

func (me *Method) Prepare() {
	// Do nothing
}

func (me *Method) Ws() {
	// Do nothing
}

func (me *Method) Ajax() {
	// Do nothing
}

func (me *Method) Get() {
	me.C.Error405()
}

func (me *Method) Post() {
	me.C.Error405()
}

func (me *Method) Delete() {
	me.C.Error405()
}

func (me *Method) Put() {
	me.C.Error405()
}

func (me *Method) Patch() {
	me.C.Error405()
}

func (me *Method) Options() {
	me.C.Error405()
}

func (me *Method) Finish() {
	// Do nothing
}

func (me *Method) getType() reflect.Type {
	me._s.RLock()
	defer me._s.RUnlock()
	return me._t
}

func (me *Method) setType(t reflect.Type) {
	me._s.Lock()
	defer me._s.Unlock()
	me._t = t
}

// Assert Serial Number
func (me *Method) asn_Core_0001() {
	// Do nothing
}

func (_ *Method) init(ro RouteHandler) {
	me := ro.(MethodInterface)
	t := me.getType()
	if t == nil {
		t = reflect.Indirect(reflect.ValueOf(me)).Type()
		me.setType(t)
	}
}

// Alais of Method
type Verb struct {
	Method
}
