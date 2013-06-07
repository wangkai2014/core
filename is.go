package core

import (
	"strings"
)

type Is struct {
	c *Core
}

func (c *Core) Is() Is {
	return Is{c}
}

// Is Ajax Request
func (i Is) AjaxRequest() bool {
	return i.c.Req.Header.Get("X-Requested-With") == "XMLHttpRequest"
}

// Is WebSocket Request
func (i Is) WebSocketRequest() bool {
	return i.c.Req.Header.Get("Connection") == "Upgrade" && i.c.Req.Header.Get("Upgrade") == "websocket"
}

// Is Do Not Track
func (i Is) DNT() bool {
	return i.c.Req.Header.Get("Dnt") == "1" || i.c.Req.Header.Get("X-Do-Not-Track") == "1"
}

// Is Secure
func (i Is) Secure() bool {
	return strings.ToLower(strings.Split(i.c.Req.Proto, "/")[0]) == "shttp"
}
