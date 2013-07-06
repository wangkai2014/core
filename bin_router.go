package core

import (
	"sort"
	"strings"
	"sync"
)

type dirLock struct {
	RouteHandler
}

func (dir dirLock) View(c *Core) {
	if c.pri.path != "" && c.pri.path != "/" {
		c.Error404()
		return
	}
	c.RouteDealer(dir.RouteHandler)
}

// Prevent Url Path Locking in BinRouter
type NoDirLock struct {
	RouteHandler
}

func (nod NoDirLock) View(c *Core) {
	c.RouteDealer(nod.RouteHandler)
}

type binRoute struct {
	dirName string
	route   RouteHandler
}

type binRoutes []*binRoute

func (bin binRoutes) Len() int {
	return len(bin)
}

func (bin binRoutes) Less(i, j int) bool {
	return bin[i].dirName < bin[j].dirName
}

func (bin binRoutes) Swap(i, j int) {
	bin[i], bin[j] = bin[j], bin[i]
}

// Binary Search Router!
type BinRouter struct {
	sync.Mutex
	routes   binRoutes
	root     RouteHandler
	asterisk RouteHandler
	group    string
}

func NewBinRouter() *BinRouter {
	return &BinRouter{routes: binRoutes{}}
}

func (bin *BinRouter) RootDir(handler RouteHandler) *BinRouter {
	bin.root = handler
	return bin
}

// Alias of RootDir
func (bin *BinRouter) Root(handler RouteHandler) *BinRouter {
	return bin.RootDir(handler)
}

func (bin *BinRouter) RootDirFunc(Func RouteHandlerFunc) *BinRouter {
	return bin.RootDir(Func)
}

// Alais of RootDirFunc
func (bin *BinRouter) RootFunc(Func RouteHandlerFunc) *BinRouter {
	return bin.RootDir(Func)
}

func (bin *BinRouter) Group(group string) *BinRouter {
	bin.group = group
	return bin
}

func (bin *BinRouter) Asterisk(handler RouteHandler) *BinRouter {
	switch t := handler.(type) {
	case *BinRouter:
		bin.asterisk = t
	case *Router:
		bin.asterisk = t
	case NoDirLock:
		bin.asterisk = t
	default:
		bin.asterisk = dirLock{handler}
	}
	return bin
}

func (bin *BinRouter) AsteriskFunc(Func RouteHandlerFunc) *BinRouter {
	return bin.Asterisk(Func)
}

func (bin *BinRouter) register(dir_ string, handler RouteHandler) {
	bin.Lock()
	defer bin.Unlock()

	dir_ = strings.TrimSpace(dir_)

	for _, route := range bin.routes {
		if route.dirName == dir_ {
			switch t := handler.(type) {
			case *BinRouter:
				route.route = t
			case *Router:
				route.route = t
			case NoDirLock:
				route.route = t
			default:
				route.route = dirLock{handler}
			}
			return
		}
	}

	switch t := handler.(type) {
	case *BinRouter:
		bin.routes = append(bin.routes, &binRoute{dir_, t})
	case *Router:
		bin.routes = append(bin.routes, &binRoute{dir_, t})
	case NoDirLock:
		bin.routes = append(bin.routes, &binRoute{dir_, t})
	default:
		bin.routes = append(bin.routes, &binRoute{dir_, dirLock{handler}})
	}
}

func (bin *BinRouter) sort() {
	bin.Lock()
	defer bin.Unlock()
	sort.Sort(bin.routes)
}

func (bin *BinRouter) Register(dir string, handler RouteHandler) *BinRouter {
	bin.register(dir, handler)
	bin.sort()
	return bin
}

func (bin *BinRouter) RegisterFunc(dir string, Func RouteHandlerFunc) *BinRouter {
	return bin.Register(dir, Func)
}

func (bin *BinRouter) RegisterMap(amap Map) *BinRouter {
	for dir, handler := range amap {
		bin.register(dir, handler)
	}
	bin.sort()
	return bin
}

func (bin *BinRouter) RegisterFuncMap(funcmap FuncMap) *BinRouter {
	for dir, handler := range funcmap {
		bin.register(dir, handler)
	}
	bin.sort()
	return bin
}

func (bin *BinRouter) error404(c *Core) {
	if !DEBUG {
		c.Error404()
		return
	}

	c.Pub.Status = 404
	out := c.Fmt()
	out.Print("404 Not Found\r\n\r\n")
	out.Print(c.Req.Host+c.pri.curpath, "\r\n\r\n")
	out.Print("Possible Directory or File!:\r\n")
	if bin.root != nil {
		out.Print("/\r\n")
	}
	for _, route := range bin.routes {
		out.Print("/", route.dirName, "\r\n")
	}
}

func (bin *BinRouter) View(c *Core) {
	if c.pri.path == "" || c.pri.path == "/" {
		if bin.root == nil {
			bin.error404(c)
			return
		}
		c.RouteDealer(bin.root)
		return
	}

	c.pri.path = strings.TrimLeft(c.pri.path, "/")
	c.pri.curpath += "/"

	pos := strings.Index(c.pri.path, "/")
	var dirname string
	if pos == -1 {
		dirname = c.pri.path
		c.pri.curpath = dirname
		c.pri.path = ""
	} else {
		dirname = c.pri.path[:pos]
		c.pri.curpath += dirname
		c.pri.path = c.pri.path[pos:]
	}

	c.Pub.BinPathDump = append(c.Pub.BinPathDump, dirname)

	if bin.group != "" {
		c.Pub.Group[bin.group] = dirname
	}

	dirname = strings.TrimSpace(dirname)

	routes_len := len(bin.routes)

	pos = sort.Search(routes_len, func(i int) bool {
		return bin.routes[i].dirName >= dirname
	})

	if pos == routes_len || bin.routes[pos].dirName != dirname {
		if bin.asterisk != nil {
			bin.asterisk.View(c)
			return
		}
		bin.error404(c)
		return
	}

	bin.routes[pos].route.View(c)
}

func SetBinRouteToMainView() {
	MainView = RouteHandlerFunc(func(c *Core) {
		appMiddlewares := AppMiddlewares.Init(c)
		defer appMiddlewares.Post()
		appMiddlewares.Pre()
		if c.CutOut() {
			return
		}

		BinRoute.View(c)
	})
}
