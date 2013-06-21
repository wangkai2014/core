package core

import (
	"sort"
	"strings"
	"sync"
)

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
	routes binRoutes
	root   RouteHandler
}

func NewBinRouter() *BinRouter {
	return &BinRouter{routes: binRoutes{}}
}

func (bin *BinRouter) RootDir(handler RouteHandler) *BinRouter {
	bin.root = handler
	return bin
}

func (bin *BinRouter) RootDirFunc(Func RouteHandlerFunc) *BinRouter {
	return bin.RootDir(Func)
}

func (bin *BinRouter) register(dir_ string, handler RouteHandler) {
	bin.Lock()
	defer bin.Unlock()

	dir_ = strings.ToLower(strings.TrimSpace(dir_))

	for _, route := range bin.routes {
		if route.dirName == dir_ {
			route.route = handler
			return
		}
	}

	bin.routes = append(bin.routes, &binRoute{dir_, handler})
}

func (bin *BinRouter) sort() {
	bin.Lock()
	defer bin.Unlock()
	sort.Sort(bin.routes)
}

func (bin *BinRouter) getRoute() binRoutes {
	bin.Lock()
	defer bin.Unlock()
	return append(binRoutes{}, bin.routes...)
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
	for _, route := range bin.getRoute() {
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
	dirname = strings.ToLower(dirname)

	routes := bin.getRoute()

	pos = sort.Search(len(routes), func(i int) bool {
		return routes[i].dirName >= dirname
	})

	if pos == len(routes) || routes[pos].dirName != dirname {
		bin.error404(c)
		return
	}

	c.RouteDealer(routes[pos].route)
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
