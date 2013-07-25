package core

import (
	"fmt"
	"regexp"
)

type Group map[string]string

func (pa Group) Set(name, value string) {
	pa[name] = value
}

func (pa Group) Get(name string) string {
	return pa[name]
}

func (pa Group) GetInt64(name string) int64 {
	num := int64(0)
	var err error
	num, err = toInt(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Group) GetInt(name string) int {
	return int(pa.GetInt64(name))
}

func (pa Group) GetInt32(name string) int32 {
	return int32(pa.GetInt64(name))
}

func (pa Group) GetInt16(name string) int16 {
	return int16(pa.GetInt64(name))
}

func (pa Group) GetInt8(name string) int8 {
	return int8(pa.GetInt64(name))
}

func (pa Group) GetUint64(name string) uint64 {
	num := uint64(0)
	var err error
	num, err = toUint(pa[name])
	if err != nil {
		return 0
	}
	return num
}

func (pa Group) GetUint(name string) uint {
	return uint(pa.GetUint64(name))
}

func (pa Group) GetUint32(name string) uint32 {
	return uint32(pa.GetUint64(name))
}

func (pa Group) GetUint16(name string) uint16 {
	return uint16(pa.GetUint64(name))
}

func (pa Group) GetUint8(name string) uint8 {
	return uint8(pa.GetUint64(name))
}

func (pa Group) GetFloat64(name string) float64 {
	num := float64(0)
	var err error
	num, err = toFloat(pa[name])
	if err != nil {
		return float64(0)
	}
	return num
}

func (pa Group) GetFloat32(name string) float32 {
	return float32(pa.GetFloat64(name))
}

type mustGroup map[string]string

func (pa mustGroup) Get(name string) string {
	return pa[name]
}

func (pa mustGroup) GetInt64(name string, c *Core) int64 {
	num := int64(0)
	var err error
	num, err = toInt(pa[name])
	if err != nil {
		c.Error404()
		return 0
	}
	return num
}

func (pa mustGroup) GetInt(name string, c *Core) int {
	return int(pa.GetInt64(name, c))
}

func (pa mustGroup) GetInt32(name string, c *Core) int32 {
	return int32(pa.GetInt64(name, c))
}

func (pa mustGroup) GetInt16(name string, c *Core) int16 {
	return int16(pa.GetInt64(name, c))
}

func (pa mustGroup) GetInt8(name string, c *Core) int8 {
	return int8(pa.GetInt64(name, c))
}

func (pa mustGroup) GetUint64(name string, c *Core) uint64 {
	num := uint64(0)
	var err error
	num, err = toUint(pa[name])
	if err != nil {
		c.Error404()
		return 0
	}
	return num
}

func (pa mustGroup) GetUint(name string, c *Core) uint {
	return uint(pa.GetUint64(name, c))
}

func (pa mustGroup) GetUint32(name string, c *Core) uint32 {
	return uint32(pa.GetUint64(name, c))
}

func (pa mustGroup) GetUint16(name string, c *Core) uint16 {
	return uint16(pa.GetUint64(name, c))
}

func (pa mustGroup) GetUint8(name string, c *Core) uint8 {
	return uint8(pa.GetUint64(name, c))
}

func (pa mustGroup) GetFloat64(name string, c *Core) float64 {
	num := float64(0)
	var err error
	num, err = toFloat(pa[name])
	if err != nil {
		c.Error404()
		return float64(0)
	}
	return num
}

func (pa mustGroup) GetFloat32(name string, c *Core) float32 {
	return float32(pa.GetFloat64(name, c))
}

type pathStr string

func (str pathStr) String() string {
	return string(str)
}

type vHostStr string

func (str vHostStr) String() string {
	return string(str)
}

type genericStr string

func (str genericStr) String() string {
	return string(str)
}

func (c *Core) pathDealer(re *regexp.Regexp, str fmt.Stringer) {
	names := re.SubexpNames()
	matches := re.FindStringSubmatch(str.String())

	for key, name := range names {
		if name != "" {
			c.Pub.Group.Set(name, matches[key])
		}
	}

	switch str.(type) {
	case pathStr:
		c.pri.curpath += matches[0]
		c.pri.path = c.pri.path[re.FindStringIndex(c.pri.path)[1]:]
	}
}
