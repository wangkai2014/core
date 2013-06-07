package core

import (
	"net/http"
)

// 'net/http.Handler' Adapter.  Implement RouterHandler interface
type HttpRouteHandler struct {
	http.Handler
}

func (ht HttpRouteHandler) View(c *Core) {
	ht.ServeHTTP(c, c.Req)
}
