package core

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/fcgi"
	"sync"
	"time"
)

// Structure of Web Application
type App struct {
	Name string

	Debug              bool
	debugTlsPortNumber uint16
	debugPortNumber    uint16

	SecureHeader string

	DefaultRouter RouteHandler
	DefaultView   RouteHandler

	TestView RouteHandler

	MiddlewareEnabled bool
	middlewares       map[string]*Middlewares
	middlewaresSync   sync.Mutex

	routers     map[string]*Router
	routersSync sync.Mutex

	dirRouters     map[string]*DirRouter
	dirRoutersSync sync.Mutex

	vHosts     map[string]*VHost
	vHostsSync sync.Mutex

	vHostsRegExp     map[string]*VHostRegExp
	vHostsRegExpSync sync.Mutex

	htmlFileCache           map[string]interface{}
	htmlFileCacheSync       sync.Mutex
	htmlGlobLocker          map[string][]string
	htmlGlobLockerSync      sync.Mutex
	HtmlTemplateCacheExpire time.Duration

	SessionCookieName          *AtomicString
	SessionExpire              time.Duration
	SessionExpireCheckInterval time.Duration
	sessionExpireCheckActive   bool
	SessionHandler             SessionHandler
	sessionMap                 map[string]sessionInterface
	sessionMapSync             sync.Mutex

	TimeLoc    *time.Location
	TimeFormat *AtomicString

	URLRev *URLReverse

	Error403 func(c *Core)
	Error404 func(c *Core)
	Error405 func(c *Core)
	Error500 func(c *Core)

	regExpCache regExpCacheSystem

	FormMemoryLimit int64

	fileServers     map[string]RouteHandler
	fileServersSync sync.Mutex

	data     map[string]interface{}
	dataSync sync.RWMutex
}

// Construct New Application
func NewApp() *App {
	app := &App{}

	app.Name = "default"

	appCount++
	if appCount > 1 {
		app.Name = fmt.Sprint(app.Name, "-", appCount)
	}

	app.middlewaresSync.Lock()
	app.routersSync.Lock()
	app.dirRoutersSync.Lock()
	app.vHostsSync.Lock()
	app.vHostsRegExpSync.Lock()
	app.htmlFileCacheSync.Lock()
	app.htmlGlobLockerSync.Lock()
	app.sessionMapSync.Lock()
	app.fileServersSync.Lock()
	app.dataSync.Lock()

	app.middlewares = map[string]*Middlewares{"main": MainMiddlewares}
	app.routers = map[string]*Router{}
	app.dirRouters = map[string]*DirRouter{}
	app.vHosts = map[string]*VHost{}
	app.vHostsRegExp = map[string]*VHostRegExp{}
	app.htmlFileCache = map[string]interface{}{}
	app.htmlGlobLocker = map[string][]string{}
	app.sessionMap = map[string]sessionInterface{}
	app.fileServers = map[string]RouteHandler{}
	app.data = map[string]interface{}{}

	app.middlewaresSync.Unlock()
	app.routersSync.Unlock()
	app.dirRoutersSync.Unlock()
	app.vHostsSync.Unlock()
	app.vHostsRegExpSync.Unlock()
	app.htmlFileCacheSync.Unlock()
	app.htmlGlobLockerSync.Unlock()
	app.sessionMapSync.Unlock()
	app.fileServersSync.Unlock()
	app.dataSync.Unlock()

	app.MiddlewareEnabled = true

	app.SessionCookieName = NewAtomicString("__session")
	app.SessionExpire = 20 * time.Minute
	app.SessionExpireCheckInterval = 10 * time.Minute
	app.SessionHandler = SessionMemory{}

	app.TimeFormat = NewAtomicString("Monday, _2 January 2006, 15:04")

	app.DefaultRouter = app.Router("main")

	app.Router("main").RegisterFunc(`^/?$`, func(c *Core) {
		c.Fmt().Print("<h1>Hello World!</h1>")
	})

	app.DefaultView = RouteHandlerFunc(func(c *Core) {
		appMiddlewares := app.Middlewares("app").Init(c)
		defer appMiddlewares.Post()
		appMiddlewares.Pre()
		if c.Terminated() {
			return
		}

		c.RouteDealer(app.DefaultRouter)
	})

	app.URLRev = &URLReverse{}

	app.Error403 = func(c *Core) {
		c.Fmt().Print("<h1>403 Forbidden</h1>")
	}
	app.Error404 = func(c *Core) {
		c.Fmt().Print("<h1>404 Not Found</h1>")
	}
	app.Error405 = func(c *Core) {
		c.Fmt().Print("<h1>405 Method Not Allowed</h1>")
	}
	app.Error500 = func(c *Core) {
		c.Fmt().Print("<h1>500 Internal Server Error</h1>")
	}

	app.regExpCache = newRegExpCacheSystem()

	app.FormMemoryLimit = 16 * 1024 * 1024

	app.SetTimeZone("Local")

	app.HtmlTemplateCacheExpire = 24 * time.Hour

	return app
}

// Get Middlewares, init on nil
func (app *App) Middlewares(name string) *Middlewares {
	app.middlewaresSync.Lock()
	defer app.middlewaresSync.Unlock()
	if app.middlewares[name] == nil {
		app.middlewares[name] = NewMiddlewares()
	}
	return app.middlewares[name]
}

// Get RegExp Router, init on nil
func (app *App) Router(name string) *Router {
	app.routersSync.Lock()
	defer app.routersSync.Unlock()
	if app.routers[name] == nil {
		app.routers[name] = NewRouter()
	}
	return app.routers[name]
}

// Get Url Directory Router, init on nil
func (app *App) DirRouter(name string) *DirRouter {
	app.dirRoutersSync.Lock()
	defer app.dirRoutersSync.Unlock()
	if app.dirRouters[name] == nil {
		app.dirRouters[name] = NewDirRouter()
	}
	return app.dirRouters[name]
}

// Get VHost, init on nil
func (app *App) VHost(name string) *VHost {
	app.vHostsSync.Lock()
	defer app.vHostsSync.Unlock()
	if app.vHosts[name] == nil {
		app.vHosts[name] = NewVHost()
	}
	return app.vHosts[name]
}

// Get VHost (Regular expression), init on nil
func (app *App) VHostRegExp(name string) *VHostRegExp {
	app.vHostsRegExpSync.Lock()
	defer app.vHostsRegExpSync.Unlock()
	if app.vHostsRegExp[name] == nil {
		app.vHostsRegExp[name] = NewVHostRegExp()
	}
	return app.vHostsRegExp[name]
}

// Get Data
func (app *App) Data(name string) interface{} {
	app.dataSync.RLock()
	defer app.dataSync.RUnlock()
	return app.data[name]
}

// Set Data
func (app *App) DataSet(name string, data interface{}) {
	app.dataSync.Lock()
	defer app.dataSync.Unlock()
	app.data[name] = data
}

// Specify File Server
func (app *App) FileServer(path, dir string) {
	app.fileServersSync.Lock()
	defer app.fileServersSync.Unlock()
	app.fileServers[path] = fileServer(path, dir)
}

// Implement http.Handler interface
func (app *App) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.serve(res, req, false)
}

func (app *App) serve(res http.ResponseWriter, req *http.Request, secure bool) {
	if app.SecureHeader != "" {
		if req.Header.Get(app.SecureHeader) != "" {
			secure = true
			req.Header.Del(app.SecureHeader)
		}
	}

	c := &Core{
		App: app,
		rw:  res.(rw),
		Req: req,
		Pub: Public{
			Status:      http.StatusOK,
			Context:     map[string]interface{}{},
			ContextStr:  map[string]string{},
			Group:       Group{},
			Session:     nil,
			TimeLoc:     app.TimeLoc,
			TimeFormat:  app.TimeFormat.String(),
			DirPathDump: []string{},
			Writers:     map[string]io.Writer{},
			Readers:     map[string]io.Reader{},
			Errors: Errors{
				E403: app.Error403,
				E404: app.Error404,
				E405: app.Error405,
				E500: app.Error500,
			},
		},
		pri: private{
			path:       req.URL.Path,
			curpath:    "",
			cut:        false,
			firstWrite: true,
			secure:     secure,
			bodyDump:   []byte{},
			allowDump:  true,
		},
	}

	c.initWriter()
	c.initTrueHost()
	c.initTrueRemoteAddr()
	c.initTruePath()
	c.initSecure()
	c.initSession()

	if app.Debug && app.TestView != nil {
		app.TestView.View(c)
		return
	}

	c.debuginfo()

	c.App.fileServersSync.Lock()
	for dir, fileServer := range c.App.fileServers {
		if len(c.pri.path) < len(dir) {
			continue
		}

		if dir == c.pri.path[:len(dir)] {
			c.App.fileServersSync.Unlock()
			fileServer.View(c)
			return
		}
	}
	c.App.fileServersSync.Unlock()

	mainMiddleware := app.Middlewares("main").Init(c)
	defer func() {
		mainMiddleware.Post()
		if !c.Terminated() && c.Req.Method != "HEAD" {
			panic(ErrorStr("No Output was sent to Client!"))
		}
	}()

	defer c.recover()

	if c.Terminated() {
		return
	}

	mainMiddleware.Pre()

	if c.Terminated() {
		return
	}

	c.RouteDealer(app.DefaultView)
}

// Start HTTP Listen
func (app *App) Listen(addr string) error {
	if app.Debug {
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			addr2 := "example.com" + addr
			_, port, _ = net.SplitHostPort(addr2)
		}
		p, _ := toUint(port)
		app.debugPortNumber = uint16(p)
	}
	mux := http.NewServeMux()
	mux.Handle("/", app)
	return http.ListenAndServe(addr, mux)
}

// Start Dummy HTTP TLS Listener
func (app *App) ListenTLSDummy(port uint16) error {
	if !app.Debug || port == 0 {
		return nil
	}
	app.debugTlsPortNumber = port
	mux := http.NewServeMux()
	mux.Handle("/", AppSecure{app})
	return http.ListenAndServe(fmt.Sprint(":", port), mux)
}

// Start HTTP TLS Listener
func (app *App) ListenTLS(addr, certFile, keyFile string) error {
	mux := http.NewServeMux()
	mux.Handle("/", AppSecure{app})
	return http.ListenAndServeTLS(addr, certFile, keyFile, mux)
}

// Start FastCGI Listener
func (app *App) ListenFCGI(l net.Listener) error {
	mux := http.NewServeMux()
	mux.Handle("/", app)
	return fcgi.Serve(l, mux)
}

// A Secure Adapter for App!
type AppSecure struct {
	*App
}

func (app AppSecure) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	app.serve(res, req, true)
}
