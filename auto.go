package core

import (
	"reflect"
	"regexp"
	"sync"
)

// Auto Populate Fields, stores expections!
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
func autoPopulateFieldByContext(c *Core, field reflect.Value, name string) {
	if c.Pub.Context[name] == nil {
		return
	}
	vcc := reflect.ValueOf(c.Pub.Context[name])
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
