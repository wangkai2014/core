package core

import (
	"net/http"
)

// Create new File Server and returns RouteHandler
func FileServer(dir string) RouteHandler {
	adir := dir
	return RouteHandlerFunc(func(c *Core) {
		c.Http().Exec(http.FileServer(http.Dir(adir)))
	})
}
