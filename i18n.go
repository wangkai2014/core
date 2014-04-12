package core

import (
	"sync"
)

type LangPackage struct {
	_m sync.RWMutex
	m  map[string]string
}

func (l *LangPackage) Key(name string) string {
	l._m.RLock()
	defer l._m.RUnlock()
	return l.m[name]
}

type Lang struct {
	_m sync.RWMutex
	m  map[string]*LangPackage
}

func (l *Lang) Package(name string) *LangPackage {
	l._m.RLock()
	defer l._m.RUnlock()
	return l.m[name]
}

func (l *Lang) Key(name string) string {
	return l.Package("core").Key(name)
}

type _langs struct {
	_m sync.RWMutex
	m  map[string]*Lang
}

var langs = &_langs{m: map[string]*Lang{}}

func (c *Context) Lang() *Lang {
	langs._m.RLock()
	defer langs._m.RUnlock()
	return langs.m[c.Pub.LangCode]
}

func LangKeyValueRegister(langCode, _package string, keyValue map[string]string) {
	langs._m.Lock()
	defer langs._m.Unlock()
	if langs.m[langCode] == nil {
		langs.m[langCode] = &Lang{m: map[string]*LangPackage{_package: &LangPackage{m: keyValue}}}
		return
	}
	langs.m[langCode]._m.Lock()
	defer langs.m[langCode]._m.Unlock()
	if langs.m[langCode].m[_package] == nil {
		langs.m[langCode].m[_package] = &LangPackage{m: keyValue}
		return
	}
	pack := langs.m[langCode].m[_package]
	pack._m.Lock()
	defer pack._m.Unlock()
	for key, value := range keyValue {
		pack.m[key] = value
	}
}

func LangAlias(aliasName, of string) {
	langs._m.Lock()
	defer langs._m.Unlock()
	if langs.m[of] == nil {
		return
	}
	langs.m[aliasName] = langs.m[of]
}

func init() {
	p := "core"

	// British English
	LangKeyValueRegister("en-GB", p, map[string]string{
		"dir":                  "ltr",
		"init":                 "initialise",
		"timeFormat":           "Monday, _2 January 2006, 15:04",
		"shortTimeFormat":      "_2/01/2006 15:04",
		"dateFormat":           "Monday, _2 January 2006",
		"shortDateFormat":      "_2/01/2006",
		"kitchenTimeFormat":    "15:04",
		"timeZoneFormat":       "MST",
		"errNoOutput":          "No output was sent to Client!",
		"err403":               "403 Forbidden",
		"err404":               "404 Not Found",
		"err405":               "405 Method Not Allowed",
		"err500":               "500 Internal Server Error",
		"errCookieNameCheck":   "Cookie name check failed",
		"errHmacDataIntegrity": "Data has been tampered with!",
	})

	// American English
	LangKeyValueRegister("en-US", p, map[string]string{
		"dir":                  "ltr",
		"init":                 "initialize",
		"timeFormat":           "Monday, January _2 2006, 3:04PM",
		"shortTimeFormat":      "01/_2/2006 3:04PM",
		"dateFormat":           "Monday, January _2 2006",
		"shortDateFormat":      "01/_2/2006",
		"kitchenTimeFormat":    "3:04PM",
		"timeZoneFormat":       "MST",
		"errNoOutput":          "No output was sent to Client!",
		"err403":               "403 Forbidden",
		"err404":               "404 Not Found",
		"err405":               "405 Method Not Allowed",
		"err500":               "500 Internal Server Error",
		"errCookieNameCheck":   "Cookie name check failed",
		"errHmacDataIntegrity": "Data has been tampered with!",
	})

	// Sadly for the British, 'en' happens to be the short version of 'en-US'
	LangAlias("en", "en-US")
}
