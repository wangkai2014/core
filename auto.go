package core

import (
	"reflect"
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
	group := mustGroup(c.Pub.Group)
	for i := 0; i < s.NumField(); i++ {
		field := s.Field(i)
		name := typeOfT.Field(i).Name
		if au.check(name) || !field.CanSet() {
			continue
		}
		if group.Get(name) == "" {
			autoPopulateFieldByContext(c, field, name)
			continue
		}
		switch field.Interface().(type) {
		case string:
			field.Set(reflect.ValueOf(group.Get(name)))
		case int:
			field.Set(reflect.ValueOf(group.GetInt(name, c)))
		case int64:
			field.Set(reflect.ValueOf(group.GetInt64(name, c)))
		case int32:
			field.Set(reflect.ValueOf(group.GetInt32(name, c)))
		case int16:
			field.Set(reflect.ValueOf(group.GetInt16(name, c)))
		case int8:
			field.Set(reflect.ValueOf(group.GetInt8(name, c)))
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
			field.Set(reflect.ValueOf(group.GetFloat32(name, c)))
		case float64:
			field.Set(reflect.ValueOf(group.GetFloat64(name, c)))
		default:
			autoPopulateFieldByContext(c, field, name)
		}
		if c.CutOut() {
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
