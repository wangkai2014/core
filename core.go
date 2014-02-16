package core

import (
	"bufio"
	"io"
	"net"
	"net/http"
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
	session    *SessionAdv
	secure     bool
	bodyDump   []byte
	allowDump  bool
}

// Strictly Public Variable
type Public struct {
	// Error Code
	Status int
	// Context, useful for storing login credentail
	Context map[string]interface{}
	// Well same as context, but for string data type only! Useful for storing user country code!
	ContextStr map[string]string
	// Used by router for storing data of named group in RegExpRule
	Group Group
	// For holding session!
	Session interface{}
	// Errors
	Errors Errors
	// Time Location
	TimeLoc *time.Location
	// Time Format
	TimeFormat string
	// DirRouter Path Dump
	DirPathDump []string
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

	if c.pri.allowDump {
		c.pri.bodyDump = append(c.pri.bodyDump, data...)
	}

	if c.pri.firstWrite {
		if c.Header().Get("Content-Type") == "" {
			c.Header().Set("Content-Type", http.DetectContentType(data))
		}

		c.pri.firstWrite = false
		c.WriteHeader(c.Pub.Status)
	}

	return c.pri.reswrite.Write(data)
}

// Response Body Dump, for Caching Purpose!
func (c *Core) BodyDump() []byte {
	return c.pri.bodyDump
}

// Disable Body Dumping
func (c *Core) NoDump() {
	c.pri.allowDump = false
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

	c.Pub.Status = num

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
func (c *Core) Terminated() bool {
	return c.pri.cut
}

// Signal framework to end user request
func (c *Core) Terminate() {
	c.pri.cut = true
}

func (c *Core) debuginfo() {
	ErrPrintf("%s, %s, %d, %s, %s, ?%s IP:%s, %v",
		c.Req.Proto, c.Req.Method, c.Pub.Status,
		c.Req.Host, c.Req.URL.Path,
		c.Req.URL.RawQuery, c.Req.RemoteAddr, time.Now())
	ErrPrintln()
}
