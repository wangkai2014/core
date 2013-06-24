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
	View(*Core)
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

func (ro *Router) load(c *Core, reset bool) bool {
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

func (ro *Router) debug(c *Core) {
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
func (ro *Router) Load(c *Core) {
	if ro.load(c, false) {
		return
	}

	if c.Is().WebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(c)
		return
	}

	c.Error404()
}

// Reset to root and try to load matching route, output 404 on fail!
func (ro *Router) LoadReset(c *Core) {
	if ro.load(c, true) {
		return
	}

	if c.Is().WebSocketRequest() {
		return
	}

	if DEBUG {
		ro.debug(c)
		return
	}

	c.Error404()
}

// Router View
func (ro *Router) View(c *Core) {
	ro.Load(c)
}

// Implement RouteHandler interface!
type RouteHandlerFunc func(*Core)

func (fn RouteHandlerFunc) View(c *Core) {
	fn(c)
}

// Reset Url, Implement RouteHandler interface!
type RouteReset struct{ *Router }

func (ro RouteReset) View(c *Core) {
	ro.LoadReset(c)
}

func (c *Core) RouteDealer(ro RouteHandler) {
	switch t := ro.(type) {
	case NoDirLock:
		t.View(c)
		return
	case MethodInterface:
		execMethodInterface(c, t)
		return
	case ProtocolInterface:
		execProtocolInterface(c, t)
		return
	case *Router:
		t.View(c)
		return
	case *BinRouter:
		t.View(c)
		return
	case RouteHandlerFunc:
		t.View(c)
		return
	case *VHost:
		t.View(c)
		return
	case *VHostRegExp:
		t.View(c)
		return
	}

	for _, routeAssert := range _getAsserter() {
		if routeAssert.Assert(c, ro) {
			return
		}
	}

	ro.View(c)
}

type RouteAsserter interface {
	Assert(*Core, RouteHandler) bool
}

type RouteAsserterFunc func(*Core, RouteHandler) bool

func (ra RouteAsserterFunc) Assert(c *Core, ro RouteHandler) bool {
	return ra(c, ro)
}

func RegisterRouteAsserter(ra ...RouteAsserter) {
	_routeAsserter.Lock()
	defer _routeAsserter.Unlock()
	_routeAsserter.ro = append(_routeAsserter.ro, ra...)
}

func RegisterRouteAsserterFunc(ra ...RouteAsserterFunc) {
	for _, raa := range ra {
		RegisterRouteAsserter(raa)
	}
}

func _getAsserter() []RouteAsserter {
	_routeAsserter.RLock()
	defer _routeAsserter.RUnlock()
	return append([]RouteAsserter{}, _routeAsserter.ro...)
}
