package core

import (
	"regexp"
	"sort"
	"sync"
)

type routerItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	Route          RouteHandler
}

// Route Handler Interface
type RouteHandler interface {
	View(*Context)
}

// Route Handler Func Map
type FuncMap map[string]RouteHandlerFunc

// Route Handler Map
type Map map[string]RouteHandler

type routes []*routerItem

func (ro routes) Len() int {
	return len(ro)
}

func (ro routes) Less(i, j int) bool {
	return ro[i].RegExp < ro[j].RegExp
}

func (ro routes) Swap(i, j int) {
	ro[i], ro[j] = ro[j], ro[i]
}

// Router (Controller), implement 'RouterHandler' interface
type Router struct {
	sync.RWMutex
	routes routes
}

func NewRouter() *Router {
	return &Router{}
}

func (ro *Router) register(RegExpRule string, handler RouteHandler) {
	ro.Lock()
	defer ro.Unlock()

	switch t := handler.(type) {
	case routeInit:
		t.init(handler)
	case RouteInit:
		t.Init(handler)
	}

	for _, route := range ro.routes {
		if route.RegExp == RegExpRule {
			route.Route = handler
			return
		}
	}

	ro.routes = append(ro.routes, &routerItem{RegExpRule, regexp.MustCompile(RegExpRule), handler})
}

func (ro *Router) sortout() {
	ro.Lock()
	defer ro.Unlock()
	sort.Sort(ro.routes)
}

// Register rule and function to Router
func (ro *Router) RegisterFunc(RegExpRule string, Function RouteHandlerFunc) *Router {
	ro.register(RegExpRule, Function)
	sort.Sort(ro.routes)
	return ro
}

// Register Map to Router, use RegExp as key!
func (ro *Router) RegisterFuncMap(funcMap FuncMap) *Router {
	if funcMap == nil {
		return ro
	}

	for rule, function := range funcMap {
		ro.register(rule, function)
	}
	ro.sortout()
	return ro
}

// Register rule and handler to Router
func (ro *Router) Register(RegExpRule string, handler RouteHandler) *Router {
	ro.register(RegExpRule, handler)
	ro.sortout()
	return ro
}

// Register Handler Map to Router, use RegExp as key!
func (ro *Router) RegisterMap(_map Map) *Router {
	if _map == nil {
		return ro
	}

	for rule, handler := range _map {
		ro.register(rule, handler)
	}
	ro.sortout()
	return ro
}

func (ro *Router) load(c *Context, reset bool) bool {
	if reset {
		c.pri.path = c.Http().Path()
		c.pri.curpath = ""
	}

	for _, route := range ro.routes {
		if !route.RegExpComplied.MatchString(c.pri.path) {
			continue
		}

		c.pathDealer(route.RegExpComplied, pathStr(c.pri.path))

		c.RouteDealer(route.Route)
		return true
	}
	return false
}

func (ro *Router) debug(c *Context) {
	c.Pub.Status = 404
	out := c.Fmt()
	out.Print("404 Not Found\r\n\r\n")
	out.Print(c.Req.Host+c.pri.curpath, "\r\n\r\n")
	out.Print("RegExp Rule(s):\r\n")
	for _, route := range ro.routes {
		out.Print(route.RegExp, "\r\n")
	}
}

// Try to load matching route, output 404 on fail!
func (ro *Router) Load(c *Context) {
	if ro.load(c, false) {
		return
	}

	if c.Is().WebSocketRequest() {
		return
	}

	if c.App.Debug {
		ro.debug(c)
		return
	}

	c.Error404()
}

// Reset to root and try to load matching route, output 404 on fail!
func (ro *Router) LoadReset(c *Context) {
	if ro.load(c, true) {
		return
	}

	if c.Is().WebSocketRequest() {
		return
	}

	if c.App.Debug {
		ro.debug(c)
		return
	}

	c.Error404()
}

// Router View
func (ro *Router) View(c *Context) {
	ro.Load(c)
}

// Implement RouteHandler interface!
type RouteHandlerFunc func(*Context)

func (fn RouteHandlerFunc) View(c *Context) {
	fn(c)
}

// Reset Url, Implement RouteHandler interface!
type RouteReset struct{ *Router }

func (ro RouteReset) View(c *Context) {
	ro.LoadReset(c)
}

type RouteAsserter interface {
	RouteHandler
	Assert(*Context, RouteHandler)
}

type routeInit interface {
	RouteHandler
	init(RouteHandler)
}

type RouteInit interface {
	RouteHandler
	Init(RouteHandler)
}

func (c *Context) RouteDealer(ro RouteHandler) {
	// Best, Average and Worst Case: O(1)
	switch t := ro.(type) {
	case MethodInterface:
		execMethodInterface(c, t)
	case ProtocolInterface:
		execProtocolInterface(c, t)
	case RouteAsserter:
		t.Assert(c, ro)
	default:
		ro.View(c)
	}
}
