package core

import (
	"reflect"
	"strings"
)

func execProtocolInterface(c *Core, pr ProtocolInterface) {
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

	switch strings.ToLower(strings.Split(c.Req.Proto, "/")[0]) {
	case "http":
		method := vc.MethodByName("Http")
		method.Call(in)
	case "shttp", "https":
		method := vc.MethodByName("Https")
		method.Call(in)
	}
}

type ProtocolInterface interface {
	View(*Core)
	Http()
	Https()
	getType() reflect.Type
	setType(reflect.Type)

	ASN_Core_0002() // Assert Serial Number
}

type Protocol struct {
	C  *Core `json:"-" xml:"-"`
	_t reflect.Type
}

func (pr *Protocol) View(c *Core) {
	pr.C = c
}

func (pr *Protocol) Http() {
	pr.C.Error404()
}

func (pr *Protocol) Https() {
	pr.C.Error404()
}

func (pr *Protocol) getType() reflect.Type {
	return pr._t
}

func (pr *Protocol) setType(t reflect.Type) {
	pr._t = t
}

// Assert Serial Number
func (pr *Protocol) ASN_Core_0002() {
	// Do nothing
}
