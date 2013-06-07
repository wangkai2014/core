package core

import (
	"reflect"
	"strings"
)

type AutoPopulateFields []string

func (au AutoPopulateFields) check(name string) bool {
	for _, item := range au {
		if name == item {
			return true
		}
	}
	return false
}

/*
Populate Struct Field Automatically
*/
func (au AutoPopulateFields) Do(c *Core, vc reflect.Value) {
	s := vc.Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		name := typeOfT.Field(i).Name
		if au.check(name) || !field.CanSet() {
			continue
		}
		if c.Pub.Group.Get(name) == "" {
			autoPopulateFieldByMeta(c, field, name)
			continue
		}
		switch field.Interface().(type) {
		case string:
			field.Set(reflect.ValueOf(c.Pub.Group.Get(name)))
		case int:
			field.Set(reflect.ValueOf(c.Pub.Group.GetInt(name)))
		case int64:
			field.Set(reflect.ValueOf(c.Pub.Group.GetInt64(name)))
		case int32:
			field.Set(reflect.ValueOf(c.Pub.Group.GetInt32(name)))
		case int16:
			field.Set(reflect.ValueOf(c.Pub.Group.GetInt16(name)))
		case int8:
			field.Set(reflect.ValueOf(c.Pub.Group.GetInt8(name)))
		case uint:
			field.Set(reflect.ValueOf(c.Pub.Group.GetUint(name)))
		case uint64:
			field.Set(reflect.ValueOf(c.Pub.Group.GetUint64(name)))
		case uint32:
			field.Set(reflect.ValueOf(c.Pub.Group.GetUint32(name)))
		case uint16:
			field.Set(reflect.ValueOf(c.Pub.Group.GetUint16(name)))
		case uint8:
			field.Set(reflect.ValueOf(c.Pub.Group.GetUint8(name)))
		case float32:
			field.Set(reflect.ValueOf(c.Pub.Group.GetFloat32(name)))
		case float64:
			field.Set(reflect.ValueOf(c.Pub.Group.GetFloat64(name)))
		default:
			autoPopulateFieldByMeta(c, field, name)
		}
	}
}

/*
Populate Field by Meta
*/
func autoPopulateFieldByMeta(c *Core, field reflect.Value, name string) {
	if c.Pub.Context[name] == nil {
		return
	}
	vcc := reflect.ValueOf(c.Pub.Context[name])
	if field.Kind() == vcc.Kind() {
		field.Set(vcc)
	}
}

func execMethodInterface(c *Core, me MethodInterface) {
	vc := reflect.New(reflect.Indirect(reflect.ValueOf(me)).Type())

	view := vc.MethodByName("View")
	in := make([]reflect.Value, 1)
	in[0] = reflect.ValueOf(c)
	view.Call(in)

	(AutoPopulateFields{"C"}).Do(c, vc)

	in = make([]reflect.Value, 0)
	method := vc.MethodByName("Prepare")
	method.Call(in)

	if c.CutOut() {
		return
	}

	is := c.Is()

	if is.WebSocketRequest() {
		method = vc.MethodByName("Ws")
		method.Call(in)
		if c.CutOut() {
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
		if c.CutOut() {
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
	View(*Core)
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

	ASN_Core_0001() // Assert Serial Number
}

type Method struct {
	C *Core
}

func (me *Method) View(c *Core) {
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

// Assert Serial Number
func (me *Method) ASN_Core_0001() {
	// Do nothing
}
