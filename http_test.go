package core

import (
	"net/http"
)

func ExampleContext_Http(c *Context) {
	// Dummy Cookie
	cookie := &http.Cookie{Name: "hello"}

	// Set Cookie
	c.Http().SetCookie(cookie)

	// Dummy Handler
	handler := func(res http.ResponseWriter, req *http.Request) {
		// Do nothing
	}

	// Execute Function
	c.Http().ExecFunc(handler)
}
