package core

import (
	"reflect"
	"strings"
)

func execProtocolInterface(c *Core, pr ProtocolInterface) {
	vc := reflect.New(reflect.Indirect(reflect.ValueOf(pr)).Type())

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

	ASN_Core_0002() // Assert Serial Number
}

type Protocol struct {
	C *Core `json:"-" xml:"-"`
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

// Assert Serial Number
func (pr *Protocol) ASN_Core_0002() {
	// Do nothing
}
