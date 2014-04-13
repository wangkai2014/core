package core

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"
)

func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// Auto Populate Fields, stores expections!
type autoPopulateFields []string

func (au autoPopulateFields) check(name string) bool {
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
func (au autoPopulateFields) do(c *Context, vc reflect.Value) {
	t := vc.Type()
	if !isStructPtr(t) {
		panic(fmt.Errorf("%v must be a struct pointer", t))
		return
	}

	s := vc.Elem()
	typeOfT := s.Type()
	group := mustGroup(c.Pub.Group)
	tag := ""

	p := func(v interface{}) interface{} {
		if tag != "positive" || c.Terminated() {
			return v
		}

		negative := false

		switch t := v.(type) {
		case int:
			negative = t < 0
		case int64:
			negative = t < 0
		case int32:
			negative = t < 0
		case int16:
			negative = t < 0
		case int8:
			negative = t < 0
		case float32:
			negative = t < 0
		case float64:
			negative = t < 0
		}

		if negative {
			c.Error404()
		}

		return v
	}

	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		name := typeOfT.Field(i).Name
		tag = typeOfT.Field(i).Tag.Get("auto")
		if au.check(name) || !field.CanSet() {
			continue
		}
		if group.Get(name) == "" {
			autoPopulateFieldByContext(c, field, name)
			continue
		}
		switch field.Interface().(type) {
		case string:
			value := group.Get(name)
			if tag == "" {
				goto value_of
			}
			if !c.App.regExpCache.Get(tag).MatchString(value) {
				c.Error404()
				return
			}
		value_of:
			field.Set(reflect.ValueOf(value))
		case int:
			field.Set(reflect.ValueOf(p(group.GetInt(name, c))))
		case int64:
			field.Set(reflect.ValueOf(p(group.GetInt64(name, c))))
		case int32:
			field.Set(reflect.ValueOf(p(group.GetInt32(name, c))))
		case int16:
			field.Set(reflect.ValueOf(p(group.GetInt16(name, c))))
		case int8:
			field.Set(reflect.ValueOf(p(group.GetInt8(name, c))))
		case uint:
			field.Set(reflect.ValueOf(group.GetUint(name, c)))
		case uint64:
			field.Set(reflect.ValueOf(group.GetUint64(name, c)))
		case uint32:
			field.Set(reflect.ValueOf(group.GetUint32(name, c)))
		case uint16:
			field.Set(reflect.ValueOf(group.GetUint16(name, c)))
		case uint8:
			field.Set(reflect.ValueOf(group.GetUint8(name, c)))
		case float32:
			field.Set(reflect.ValueOf(p(group.GetFloat32(name, c))))
		case float64:
			field.Set(reflect.ValueOf(p(group.GetFloat64(name, c))))
		default:
			autoPopulateFieldByContext(c, field, name)
		}
		if c.Terminated() {
			return
		}
	}
}

/*
Populate Field by Meta
*/
func autoPopulateFieldByContext(c *Context, field reflect.Value, name string) {
	if c.Pub.Data[name] == nil {
		return
	}
	vcc := reflect.ValueOf(c.Pub.Data[name])
	if field.Kind() == vcc.Kind() {
		field.Set(vcc)
	}
}

type regExpCacheSystem struct {
	sync.Mutex
	res map[string]*regexp.Regexp
}

func newRegExpCacheSystem() regExpCacheSystem {
	return regExpCacheSystem{res: map[string]*regexp.Regexp{}}
}

func (reg regExpCacheSystem) Get(str string) *regexp.Regexp {
	reg.Lock()
	defer reg.Unlock()
	if reg.res[str] != nil {
		return reg.res[str]
	}
	_re := regexp.MustCompile(str)
	reg.res[str] = _re
	return _re
}

type Auto struct {
	c *Context
}

func (c *Context) Auto() Auto {
	return Auto{c}
}

// Auto Populate from c.Pub.Group and c.Pub.Data
func (a Auto) PopulateStructFieldsValue(structPointer reflect.Value, exclude ...string) {
	(autoPopulateFields(exclude)).do(a.c, structPointer)
}

// Auto Populate from c.Pub.Group and c.Pub.Data
func (a Auto) PopulateStructFields(structPointer interface{}, exclude ...string) {
	a.PopulateStructFieldsValue(reflect.ValueOf(structPointer), exclude...)
}
