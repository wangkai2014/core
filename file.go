package core

import (
	"net/http"
)

// Create new File Server and returns RouteHandler
func FileServer(dir string) RouteHandler {
	return NoDirLock{RouteHandlerFunc(func(c *Context) {
		c.Http().Exec(http.FileServer(http.Dir(dir)))
	})}
}

func fileServer(path, dir string) RouteHandler {
	return RouteHandlerFunc(func(c *Context) {
		http.StripPrefix(path, http.FileServer(http.Dir(dir))).ServeHTTP(c.Res, c.Req)
	})
}
