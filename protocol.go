package core

import (
	"reflect"
	"sync"
)

func execProtocolInterface(c *Context, pr ProtocolInterface) {
	t := pr.getType()
	if t == nil {
		t = reflect.Indirect(reflect.ValueOf(pr)).Type()
		pr.setType(t)
	}

	vc := reflect.New(t)

	view := vc.MethodByName("View")
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(c)
	view.Call(in)

	in = make([]reflect.Value, 0)

	switch c.pri.secure {
	case false:
		method := vc.MethodByName("Http")
		method.Call(in)
	case true:
		method := vc.MethodByName("Https")
		method.Call(in)
	}
}

type ProtocolInterface interface {
	View(*Context)
	Http()
	Https()
	getType() reflect.Type
	setType(reflect.Type)

	asn_Core_0002() // Assert Serial Number
}

type Protocol struct {
	C  *Context `json:"-" xml:"-"`
	_t reflect.Type
	_s sync.RWMutex
}

func (pr *Protocol) View(c *Context) {
	pr.C = c
}

func (pr *Protocol) Http() {
	pr.C.Error404()
}

func (pr *Protocol) Https() {
	pr.C.Error404()
}

func (pr *Protocol) getType() reflect.Type {
	pr._s.RLock()
	defer pr._s.RUnlock()
	return pr._t
}

func (pr *Protocol) setType(t reflect.Type) {
	pr._s.Lock()
	defer pr._s.Unlock()
	pr._t = t
}

// Assert Serial Number
func (pr *Protocol) asn_Core_0002() {
	// Do nothing
}

func (_ *Protocol) init(ro RouteHandler) {
	pr := ro.(ProtocolInterface)
	t := pr.getType()
	if t == nil {
		t = reflect.Indirect(reflect.ValueOf(pr)).Type()
		pr.setType(t)
	}
}
