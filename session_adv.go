package core

import (
	"encoding/gob"
)

type sessionAdvMap map[string]interface{}

func init() {
	gob.Register(sessionAdvMap{})
}

// Key/Value version of Session
type SessionAdv struct {
	se   Session
	sMap sessionAdvMap
}

// Key/Value version of Session
func (se Session) Adv() *SessionAdv {
	if se.c.pri.session != nil {
		return se.c.pri.session
	}
	se.c.pri.session = &SessionAdv{se: se}
	se.c.pri.session.init()
	return se.c.pri.session
}

func (se *SessionAdv) init() {
	if se.sMap != nil {
		return
	}

	switch t := se.se.c.Pub.Session.(type) {
	case sessionAdvMap:
		se.sMap = t
	default:
		se.sMap = sessionAdvMap{}
	}
}

// Set Session by Key
func (se *SessionAdv) Set(key string, value interface{}) {
	se.sMap[key] = value
}

// Get Session by Key
func (se *SessionAdv) Get(key string) interface{} {
	return se.sMap[key]
}

// Save Session and set Cookie to client.
func (se *SessionAdv) Save() {
	se.se.Set(se.sMap)
}
