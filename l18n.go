package core

import (
	"sync"
)

type LangPackage struct {
	sync.RWMutex
	m map[string]string
}

func (l *LangPackage) Key(name string) string {
	l.RLock()
	defer l.RUnlock()
	return l.m[name]
}

type Lang struct {
	sync.RWMutex
	m map[string]*LangPackage
}

func (l *Lang) Package(name string) *LangPackage {
	l.RLock()
	defer l.RUnlock()
	return l.m[name]
}

func (l *Lang) Key(name string) string {
	return l.Package("core").Key(name)
}

type _langs struct {
	sync.RWMutex
	m map[string]*Lang
}

var langs = &_langs{m: map[string]*Lang{}}

func (c *Context) Lang() *Lang {
	langs.RLock()
	defer langs.RUnlock()
	return langs.m[c.Pub.LangCode]
}

func LangKeyValueRegister(langCode, _package string, keyValue map[string]string) {
	langs.Lock()
	defer langs.Unlock()
	if langs.m[langCode] == nil {
		langs.m[langCode] = &Lang{m: map[string]*LangPackage{_package: &LangPackage{m: keyValue}}}
		return
	}
	langs.m[langCode].Lock()
	defer langs.m[langCode].Unlock()
	if langs.m[langCode].m[_package] == nil {
		langs.m[langCode].m[_package] = &LangPackage{m: keyValue}
		return
	}
	pack := langs.m[langCode].m[_package]
	pack.Lock()
	defer pack.Unlock()
	for key, value := range keyValue {
		pack.m[key] = value
	}
}

func LangAlias(aliasName, of string) {
	langs.Lock()
	defer langs.Unlock()
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
