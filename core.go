package core

import (
	"bufio"
	"fmt"
	html "html/template"
	"io"
	"net"
	"net/http"
	"net/http/cgi"
	"net/http/fcgi"
	"os"
	"runtime"
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
}

type Public struct {
	// Error Code
	Status int
	// Meta, useful for storing login credentail
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
}

// The Framework Structure, it's implement the interfaces of 'net/http.ResponseWriter',
// 'net/http.Hijacker', 'net/http.Flusher' and 'net/http.Handler'
type Core struct {
	// Request
	Req *http.Request
	// Public Variabless
	Pub Public
	rw
	pri private
}

// HTTP Handler
func (_ Core) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	c := &Core{
		rw:  res.(rw),
		Req: req,
		Pub: Public{
			Status:      http.StatusOK,
			Context:     map[string]interface{}{},
			ContextStr:  map[string]string{},
			Group:       Group{},
			HtmlFunc:    html.FuncMap{},
			Session:     nil,
			TimeLoc:     DefaultTimeLoc,
			TimeFormat:  DefaultTimeFormat.String(),
			BinPathDump: []string{},
			Errors: Errors{
				E403: Error403,
				E404: Error404,
				E405: Error405,
				E500: Error500,
			},
		},
		pri: private{
			path:       req.URL.Path,
			curpath:    "",
			cut:        false,
			firstWrite: true,
		},
	}

	c.initWriter()
	c.initTrueHost()
	c.initTrueRemoteAddr()
	c.initTruePath()
	c.initSecure()
	c.initSession()

	defer c.recover()

	mainMiddleware := MainMiddlewares.Init(c)
	defer func() {
		defer c.recover()
		mainMiddleware.Post()
		if !c.CutOut() {
			panic(ErrorStr("No Output was sent to Client!"))
		}
	}()

	c.debugStart()
	defer c.debugEnd()

	mainMiddleware.Html()

	if c.CutOut() {
		return
	}

	mainMiddleware.Pre()

	if c.CutOut() {
		return
	}

	c.RouteDealer(MainView)

	if c.CutOut() {
		return
	}
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

func (c *Core) debuginfo(a string) {
	if !DEBUG {
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

var _core = Core{}

func secure(res http.ResponseWriter, req *http.Request) {
	req.Header.Set("X-Secure-Mode", "1")
	_core.ServeHTTP(res, req)
}

func nonsecure(res http.ResponseWriter, req *http.Request) {
	req.Header.Del("X-Secure-Mode")
	_core.ServeHTTP(res, req)
}

// Start Http Server
func StartHttp(addr string) error {
	return http.ListenAndServe(addr, http.HandlerFunc(nonsecure))
}

// Start Http Server with TLS
func StartHttpTLS(addr string, certFile string, keyFile string) error {
	return http.ListenAndServeTLS(addr, certFile, keyFile, http.HandlerFunc(secure))
}

// Emulate TLS Server via Http connection, the connection is unencrypted!
//
// Only work in debug mode.  Does not require certificate files.
func StartDummyHttpTLS(port uint16) error {
	if !DEBUG {
		return nil
	}
	debugTlsPortNumber = port
	return http.ListenAndServe(fmt.Sprint(":", port), http.HandlerFunc(secure))
}

// Start FastCGI Server
func StartFastCGI(l net.Listener) error {
	if l == nil {
		os.Stderr = nil
	}
	return fcgi.Serve(l, http.HandlerFunc(nonsecure))
}

// Start CGI, disables Stderr completely. (Due to the way how IIS handlers Stderr)
func StartCGI() error {
	os.Stderr = nil
	return cgi.Serve(http.HandlerFunc(nonsecure))
}

type DebugMiddleware struct {
	Middleware
}

func (de *DebugMiddleware) Html() {
	c := de.C
	c.Pub.HtmlFunc["Debug"] = func() bool {
		return DEBUG
	}

	c.Pub.HtmlFunc["NotDebug"] = func() bool {
		return !DEBUG
	}
}

func init() {
	MainMiddlewares.Register(&DebugMiddleware{})
}

// Use Max CPU
func UseMaxCPU() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
