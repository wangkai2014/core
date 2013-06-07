package core

import (
	"net/http"
)

func ExampleCore_Http(c *Core) {
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
