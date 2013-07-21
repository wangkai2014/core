package core

import (
	"bufio"
	"fmt"
	html "html/template"
	"io"
	"net"
	"net/http"
	"net/http/fcgi"
	//"os"
	"sync"
	"time"
)

type rw interface {
	http.ResponseWriter
}

type private struct {
	path       string
	curpath    string
	reswrite   io.Writer
	cut        bool
	firstWrite bool
	html       *htmlDefault
	session    *SessionAdv
	secure     bool
}

type Public struct {
	// Error Code
	Status int
	// Context, useful for storing login credentail
	Context map[string]interface{}
	// Well same as context, but for string data type only! Useful for storing user country code!
	ContextStr map[string]string
	// Used by router for storing data of named group in RegExpRule
	Group Group
	// Function to load in html template system.
	HtmlFunc html.FuncMap
	// For holding session!
	Session interface{}
	// Errors
	Errors Errors
	// Time Location
	TimeLoc *time.Location
	// Time Format
	TimeFormat string
	// BinRouter Path Dump
	BinPathDump []string
	// Reader and Writer Dump
	Writers map[string]io.Writer
	Readers map[string]io.Reader
}

// The Framework Structure, it's implement the interfaces of 'net/http.ResponseWriter',
// 'net/http.Hijacker', 'net/http.Flusher' and 'net/http.Handler'
type Core struct {
	App *App
	// Request
	Req *http.Request
	// Public Variabless
	Pub Public
	rw
	pri private
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (c *Core) Header() http.Header {
	return c.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (c *Core) Write(data []byte) (int, error) {
	c.pri.cut = true

	if c.pri.firstWrite {
		if c.Header().Get("Content-Type") == "" {
			c.Header().Set("Content-Type", http.DetectContentType(data))
		}

		c.pri.firstWrite = false
		c.WriteHeader(c.Pub.Status)
	}

	return c.pri.reswrite.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
//
// Note: Use Status property to set error code! As this disable compression!
func (c *Core) WriteHeader(num int) {
	c.pri.cut = true

	if c.pri.firstWrite {
		c.pri.firstWrite = false
	}

	c.rw.WriteHeader(num)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the HTTP server library
// will not do anything else with the connection.
// It becomes the caller's responsibility to manage
// and close the connection.
func (c *Core) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	c.pri.cut = true

	hj, ok := c.rw.(http.Hijacker)
	if ok {
		return hj.Hijack()
	}

	return nil, nil, ErrorStr("Connection is not Hijackable")
}

// Flush sends any buffered data to the client.
func (c *Core) Flush() {
	fl, ok := c.rw.(http.Flusher)
	if ok {
		fl.Flush()
	}
}

// true if output was sent to client, otherwise false!
func (c *Core) CutOut() bool {
	return c.pri.cut
}

func (c *Core) Cut() {
	c.pri.cut = true
}

func (c *Core) debuginfo(a string) {
	if !c.App.Debug {
		return
	}
	ErrPrintf("--\r\n %s  %s, %s, %s, %s, ?%s IP:%s \r\n--\r\n",
		a, c.Req.Proto, c.Req.Method,
		c.Req.Host, c.Req.URL.Path,
		c.Req.URL.RawQuery, c.Req.RemoteAddr)
}

func (c *Core) debugStart() {
	c.debuginfo("START")
}

func (c *Core) debugEnd() {
	c.debuginfo("END  ")
}

type App struct {
	Name string

	Debug              bool
	debugTlsPortNumber uint16

	SecureHeader string

	DefaultRouter RouteHandler
	DefaultView   RouteHandler

	MiddlewareEnabled bool
	middlewares       map[string]*Middlewares
	middlewaresSync   sync.Mutex

	routers     map[string]*Router
	routersSync sync.Mutex

	binRouters     map[string]*BinRouter
	binRoutersSync sync.Mutex

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

func NewApp() *App {
	app := &App{}

	app.middlewaresSync.Lock()
	app.routersSync.Lock()
	app.binRoutersSync.Lock()
	app.vHostsSync.Lock()
	app.vHostsRegExpSync.Lock()
	app.htmlFileCacheSync.Lock()
	app.htmlGlobLockerSync.Lock()
	app.sessionMapSync.Lock()
	app.fileServersSync.Lock()
	app.dataSync.Lock()

	app.middlewares = map[string]*Middlewares{"main": MainMiddlewares}
	app.routers = map[string]*Router{}
	app.binRouters = map[string]*BinRouter{}
	app.vHosts = map[string]*VHost{}
	app.vHostsRegExp = map[string]*VHostRegExp{}
	app.htmlFileCache = map[string]interface{}{}
	app.htmlGlobLocker = map[string][]string{}
	app.sessionMap = map[string]sessionInterface{}
	app.fileServers = map[string]RouteHandler{}
	app.data = map[string]interface{}{}

	app.middlewaresSync.Unlock()
	app.routersSync.Unlock()
	app.binRoutersSync.Unlock()
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
		if c.CutOut() {
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

func (app *App) Middlewares(name string) *Middlewares {
	app.middlewaresSync.Lock()
	defer app.middlewaresSync.Unlock()
	if app.middlewares[name] == nil {
		app.middlewares[name] = NewMiddlewares()
	}
	return app.middlewares[name]
}

func (app *App) Router(name string) *Router {
	app.routersSync.Lock()
	defer app.routersSync.Unlock()
	if app.routers[name] == nil {
		app.routers[name] = NewRouter()
	}
	return app.routers[name]
}

func (app *App) BinRouter(name string) *BinRouter {
	app.binRoutersSync.Lock()
	defer app.binRoutersSync.Unlock()
	if app.binRouters[name] == nil {
		app.binRouters[name] = NewBinRouter()
	}
	return app.binRouters[name]
}

func (app *App) VHost(name string) *VHost {
	app.vHostsSync.Lock()
	defer app.vHostsSync.Unlock()
	if app.vHosts[name] == nil {
		app.vHosts[name] = NewVHost()
	}
	return app.vHosts[name]
}

func (app *App) VHostRegExp(name string) *VHostRegExp {
	app.vHostsRegExpSync.Lock()
	defer app.vHostsRegExpSync.Unlock()
	if app.vHostsRegExp[name] == nil {
		app.vHostsRegExp[name] = NewVHostRegExp()
	}
	return app.vHostsRegExp[name]
}

func (app *App) Data(name string) interface{} {
	app.dataSync.RLock()
	defer app.dataSync.RUnlock()
	return app.data[name]
}

func (app *App) DataSet(name string, data interface{}) {
	app.dataSync.Lock()
	defer app.dataSync.Unlock()
	app.data[name] = data
}

func (app *App) FileServer(path, dir string) {
	app.fileServersSync.Lock()
	defer app.fileServersSync.Unlock()
	app.fileServers[path] = fileServer(path, dir)
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
			HtmlFunc:    html.FuncMap{},
			Session:     nil,
			TimeLoc:     app.TimeLoc,
			TimeFormat:  app.TimeFormat.String(),
			BinPathDump: []string{},
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
		},
	}

	c.initWriter()
	c.initTrueHost()
	c.initTrueRemoteAddr()
	c.initTruePath()
	c.initSecure()
	c.initSession()

	c.debugStart()
	defer c.debugEnd()

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
		if !c.CutOut() && c.Req.Method != "HEAD" {
			panic(ErrorStr("No Output was sent to Client!"))
		}
	}()

	defer c.recover()

	mainMiddleware.Html()

	if c.CutOut() {
		return
	}

	mainMiddleware.Pre()

	if c.CutOut() {
		return
	}

	c.RouteDealer(app.DefaultView)
}

func (app *App) Listen(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.serve(res, req, false)
	}))
}

func (app *App) ListenTLSDummy(port uint16) error {
	if !app.Debug || port == 0 {
		return nil
	}
	app.debugTlsPortNumber = port
	return http.ListenAndServe(fmt.Sprint(":", port),
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			app.serve(res, req, true)
		}),
	)
}

func (app *App) ListenTLS(addr, certFile, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile,
		http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			app.serve(res, req, true)
		}),
	)
}

func (app *App) ListenFCGI(l net.Listener) error {
	return fcgi.Serve(l, http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		app.serve(res, req, false)
	}))
}

// Debug Middleware, Add HTML Function to template!
type DebugMiddleware struct {
	Middleware
}

func (de *DebugMiddleware) Html() {
	c := de.C
	c.Pub.HtmlFunc["Debug"] = func() bool {
		return de.C.App.Debug
	}

	c.Pub.HtmlFunc["NotDebug"] = func() bool {
		return !de.C.App.Debug
	}
}

func init() {
	MainMiddlewares.Register(&DebugMiddleware{})
}
