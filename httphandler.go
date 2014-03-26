package core

import (
	"net/http"
)

// 'net/http.Handler' Adapter.  Implement RouterHandler interface
type HttpRouteHandler struct {
	http.Handler
}

// Implement RouteHandler
func (ht HttpRouteHandler) View(c *Context) {
	ht.ServeHTTP(c.Res, c.Req)
}
