package core

import (
	"sync"
	"time"
)

var (
	/*
		Debug
	*/
	DEBUG              = false
	debugTlsPortNumber = uint16(0)

	/*
		Main View
	*/
	MainView RouteHandler = RouteHandlerFunc(func(c *Core) {
		appMiddlewares := AppMiddlewares.Init(c)
		defer appMiddlewares.Post()
		appMiddlewares.Pre()
		if c.CutOut() {
			return
		}

		Route.Load(c)
	})

	/*
		Middleware
	*/
	MainMiddlewares = NewMiddlewares()
	AppMiddlewares  = NewMiddlewares()

	/*
		Route
	*/
	Route    = NewRouter()
	BinRoute = NewBinRouter()

	/*
		Form
	*/
	FormMemoryLimit = int64(16 * 1024 * 1024)

	/*
		Html
	*/
	htmlFileCache = struct {
		sync.Mutex
		m map[string]interface{}
	}{m: map[string]interface{}{}}
	HtmlTemplateCacheExpire = 24 * time.Hour
	htmlGlobLocker          = struct {
		sync.Mutex
		filenames map[string][]string
	}{
		filenames: map[string][]string{},
	}

	/*
		Session Settings
	*/
	SessionCookieName                         = NewAtomicString("__session") // Session Cookie Name
	SessionExpire              time.Duration  = 20 * time.Minute             // Session Expiry
	SessionExpiryCheckInterval time.Duration  = 10 * time.Minute             // Session Expire Check Interval
	sessionExpiryCheckActive                  = false
	DefaultSessionHandler      SessionHandler = SessionMemory{} // Default Session Handler

	/*
		Time
	*/
	DefaultTimeLoc    *time.Location                                      // Default Time Location
	DefaultTimeFormat = NewAtomicString("Monday, _2 January 2006, 15:04") // Default Time Format

	/*
		Url Reverse
	*/
	URLRev = &URLReverse{}

	/*
		Error Handler
	*/
	Error403 = func(c *Core) {
		c.Fmt().Print("<h1>403 Forbidden</h1>")
	}
	Error404 = func(c *Core) {
		c.Fmt().Print("<h1>404 Not Found</h1>")
	}
	Error405 = func(c *Core) {
		c.Fmt().Print("<h1>405 Method Not Allowed</h1>")
	}
	Error500 = func(c *Core) {
		c.Fmt().Print("<h1>500 Internal Server Error</h1>")
	}

	/*
		VHost Vars
	*/
	VHosts       = NewVHost()
	VHostsRegExp = NewVHostRegExp()

	_routeAsserter = struct {
		sync.RWMutex
		ro []RouteAsserter
	}{ro: []RouteAsserter{}}
)
