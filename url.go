package core

import (
	"fmt"
	"strings"
	"sync"
)

// URL Reverse Map
type URLReverseMap map[string]string

// URL Reverse Data Type
type URLReverse struct {
	sync.RWMutex
	urls URLReverseMap
}

func (u *URLReverse) register(name, format string) {
	u.Lock()
	defer u.Unlock()
	u.urls[name] = format
}

func (u *URLReverse) initUrls() {
	u.Lock()
	defer u.Unlock()
	if u.urls == nil {
		u.urls = URLReverseMap{}
	}
}

// Url Reverse Register
func (u *URLReverse) Register(name, format string) *URLReverse {
	u.initUrls()
	u.register(name, format)
	return u
}

// Url Reverse Register Map
func (u *URLReverse) RegisterMap(urls URLReverseMap) *URLReverse {
	if urls == nil {
		return u
	}

	u.initUrls()

	for name, format := range urls {
		u.register(name, format)
	}

	return u
}

// Print relative URL to string
func (u *URLReverse) Print(name string, a ...interface{}) string {
	u.RLock()
	defer u.RUnlock()
	return fmt.Sprintf(u.urls[name], a...)
}

type Url struct {
	c *Core
}

func (c *Core) Url() Url {
	return Url{c}
}

// Get Absolute URL, you can leave relative_url blank just to get the root url.
func (u Url) Absolute(relative_url string) string {
	if u.c.Is().Secure() {
		return u.AbsoluteHttps(relative_url)
	}
	return u.AbsoluteHttp(relative_url)
}

// Get Absolute URL (http), you can leave relative_url blank just to get the root url.
func (u Url) AbsoluteHttp(relative_url string) string {
	c := u.c
	if c.Req.URL.Host != "" {
		return "http://" + c.Req.URL.Host + relative_url
	}

	return relative_url
}

// Get Absolute URL (https), you can leave relative_url blank just to get the root url.
func (u Url) AbsoluteHttps(relative_url string) string {
	c := u.c
	if c.Req.URL.Host != "" {
		host := c.Req.URL.Host
		protocol := "https://"
		if debugTlsPortNumber != uint16(0) {
			protocol = "http://"
			host = strings.Split(host, ":")[0] + fmt.Sprint(":", debugTlsPortNumber)
		}
		return protocol + host + relative_url
	}

	return relative_url
}

// Shortcut to Reverse
func (u Url) Reverse(name string, a ...interface{}) string {
	return URLRev.Print(name, a...)
}

func (u Url) code301() int {
	if u.c.Pub.Status == 200 {
		return 301
	}
	return u.c.Pub.Status
}

func (u Url) code303() int {
	if u.c.Pub.Status == 200 {
		return 303
	}
	return u.c.Pub.Status
}

// Convert current path to Https
func (u Url) ToHttps() {
	defer u.c.WriteHeader(u.code301())
	u.c.Header().Set("Location", u.AbsoluteHttps(u.c.Req.URL.Path))
}

// Convert current path to Http
func (u Url) ToHttp() {
	defer u.c.WriteHeader(u.code301())
	u.c.Header().Set("Location", u.AbsoluteHttp(u.c.Req.URL.Path))
}

// Redirect client to relative_url
func (u Url) Redirect(relative_url string) {
	defer u.c.WriteHeader(u.code303())
	u.c.Header().Set("Location", u.Absolute(relative_url))
}

// Redirect client to relative_url (Http Only)
func (u Url) RedirectHttp(relative_url string) {
	defer u.c.WriteHeader(u.code303())
	u.c.Header().Set("Location", u.AbsoluteHttp(relative_url))
}

// Redirect client to relative_url (Https Only)
func (u Url) RedirectHttps(relative_url string) {
	defer u.c.WriteHeader(u.code303())
	u.c.Header().Set("Location", u.AbsoluteHttps(relative_url))
}

type URLReverseMiddleware struct {
	Middleware
}

func (url *URLReverseMiddleware) Html() {
	url.C.Pub.HtmlFunc["url"] = func(name string, a ...interface{}) string {
		return URLRev.Print(name, a...)
	}
}

func init() {
	MainMiddlewares.Register(&URLReverseMiddleware{})
}
