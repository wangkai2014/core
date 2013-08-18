package core

import (
	"net/http"
)

// 'net/http.Handler' Adapter.  Implement RouterHandler interface
type HttpRouteHandler struct {
	http.Handler
}

// Implement RouteHandler
func (ht HttpRouteHandler) View(c *Core) {
	ht.ServeHTTP(c, c.Req)
}
