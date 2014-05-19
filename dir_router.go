package core

import (
	"regexp"
	"strings"
	"sync"
)

type dirLock struct {
	RouteHandler
}

func (dir dirLock) View(c *Context) {
	if c.pri.path != "" && c.pri.path != "/" {
		c.Error404()
		return
	}
	c.RouteDealer(dir.RouteHandler)
}

// Prevent Url Path Locking in DirRouter
type NoDirLock struct {
	RouteHandler
}

// Implement RouteHandler
func (nod NoDirLock) View(c *Context) {
	c.RouteDealer(nod.RouteHandler)
}

type dirRoute struct {
	dirName string
	route   RouteHandler
}

// Directory Search Router!
type DirRouter struct {
	sync.RWMutex
	routes   map[string]*dirRoute
	root     RouteHandler
	asterisk RouteHandler
	group    string
	regexp   *regexp.Regexp
}

// Construct Directory Router
func NewDirRouter() *DirRouter {
	return &DirRouter{routes: map[string]*dirRoute{}}
}

// Set Root Directory Handler
func (dir *DirRouter) RootDir(handler RouteHandler) *DirRouter {
	dir.root = handler
	return dir
}

// Alias of RootDir
func (dir *DirRouter) Root(handler RouteHandler) *DirRouter {
	return dir.RootDir(handler)
}

// Set Root Directory Function
func (dir *DirRouter) RootDirFunc(Func RouteHandlerFunc) *DirRouter {
	return dir.RootDir(Func)
}

// Alais of RootDirFunc
func (dir *DirRouter) RootFunc(Func RouteHandlerFunc) *DirRouter {
	return dir.RootDir(Func)
}

// Set Group Name
func (dir *DirRouter) Group(group string) *DirRouter {
	dir.group = group
	return dir
}

// Set Asterisk Handler
func (dir *DirRouter) Asterisk(handler RouteHandler) *DirRouter {
	switch t := handler.(type) {
	case *DirRouter:
		dir.asterisk = t
	case *Router:
		dir.asterisk = t
	case NoDirLock:
		dir.asterisk = t
	default:
		dir.asterisk = dirLock{handler}
	}
	return dir
}

// Set Asterisk Function
func (dir *DirRouter) AsteriskFunc(Func RouteHandlerFunc) *DirRouter {
	return dir.Asterisk(Func)
}

// Set Regular Expression
func (dir *DirRouter) RegExp(pattern string) *DirRouter {
	dir.regexp = regexp.MustCompile(pattern)
	return dir
}

func (dir *DirRouter) register(dir_ string, handler RouteHandler) {
	dir.Lock()
	defer dir.Unlock()

	if strings.ContainsAny(dir_, `/\`) {
		return
	}

	switch t := handler.(type) {
	case routeInit:
		t.init(handler)
	case RouteInit:
		t.Init(handler)
	}

	dir_ = strings.TrimSpace(dir_)

	switch t := handler.(type) {
	case *DirRouter:
		dir.routes[dir_] = &dirRoute{dir_, t}
	case *Router:
		dir.routes[dir_] = &dirRoute{dir_, t}
	case NoDirLock:
		dir.routes[dir_] = &dirRoute{dir_, t}
	default:
		dir.routes[dir_] = &dirRoute{dir_, dirLock{handler}}
	}
}

// Register Handler to Directory
func (dir *DirRouter) Register(dir_ string, handler RouteHandler) *DirRouter {
	dir.register(dir_, handler)
	return dir
}

// Register Function to Directory
func (dir *DirRouter) RegisterFunc(dir_ string, Func RouteHandlerFunc) *DirRouter {
	return dir.Register(dir_, Func)
}

// Register Map of Handler ("dir": handler)
func (dir *DirRouter) RegisterMap(amap Map) *DirRouter {
	for dir_, handler := range amap {
		dir.register(dir_, handler)
	}
	return dir
}

// Register Map of Functions ("dir": function)
func (dir *DirRouter) RegisterFuncMap(funcmap FuncMap) *DirRouter {
	for dir_, handler := range funcmap {
		dir.register(dir_, handler)
	}
	return dir
}

func (dir *DirRouter) error404(c *Context) {
	if !c.App.Debug {
		c.Error404()
		return
	}

	c.Pub.Status = 404
	out := c.Fmt()
	out.Print("404 Not Found\r\n\r\n")
	out.Print(c.Req.Host+c.pri.curpath, "\r\n\r\n")
	out.Print("Possible Directory or File!:\r\n")
	if dir.root != nil {
		out.Print("/\r\n")
	}
	for _, route := range dir.routes {
		out.Print("/", route.dirName, "\r\n")
	}
}

// Implement RouteHandler
func (dir *DirRouter) View(c *Context) {
	// Check if Root Path
	if c.pri.path == "" || c.pri.path == "/" {
		if dir.root == nil {
			dir.error404(c)
			return
		}
		c.RouteDealer(dir.root)
		return
	}

	c.pri.path = strings.TrimLeft(c.pri.path, "/")
	c.pri.curpath += "/"

	pos := strings.Index(c.pri.path, "/")
	var dirname string
	if pos == -1 {
		dirname = c.pri.path
		c.pri.curpath += dirname
		c.pri.path = ""
	} else {
		dirname = c.pri.path[:pos]
		c.pri.curpath += dirname
		c.pri.path = c.pri.path[pos:]
	}

	c.Pub.DirPathDump = append(c.Pub.DirPathDump, dirname)

	dirname = strings.TrimSpace(dirname)

	dir.RLock()
	route := dir.routes[dirname]
	dir.RUnlock()

	if route == nil {
		if dir.asterisk != nil {
			if dir.regexp != nil {
				if !dir.regexp.MatchString(dirname) {
					dir.error404(c)
					return
				}
				c.pathDealer(dir.regexp, genericStr(dirname))
			} else if dir.group != "" {
				c.Pub.Group[dir.group] = dirname
			}
			dir.asterisk.View(c)
			return
		}
		dir.error404(c)
		return
	}

	route.route.View(c)
}
