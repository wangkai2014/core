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

type Res struct {
	rw
	c *Context
}

// Header returns the header map that will be sent by WriteHeader.
// Changing the header after a call to WriteHeader (or Write) has
// no effect.
func (r Res) Header() http.Header {
	return r.rw.Header()
}

// Write writes the data to the connection as part of an HTTP reply.
// If WriteHeader has not yet been called, Write calls WriteHeader(http.StatusOK)
// before writing the data.  If the Header does not contain a
// Content-Type line, Write adds a Content-Type set to the result of passing
// the initial 512 bytes of written data to DetectContentType.
func (r Res) Write(data []byte) (int, error) {
	r.c.pri.cut = true

	if r.c.pri.firstWrite {
		if r.Header().Get("Content-Type") == "" {
			r.Header().Set("Content-Type", http.DetectContentType(data))
		}

		r.c.pri.firstWrite = false
		r.WriteHeader(r.c.Pub.Status)
	}

	return r.c.pri.reswrite.Write(data)
}

// WriteHeader sends an HTTP response header with status code.
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
func (r Res) WriteHeader(num int) {
	r.c.pri.cut = true

	if r.c.pri.firstWrite {
		r.c.pri.firstWrite = false
	}

	r.c.Pub.Status = num

	r.rw.WriteHeader(num)
}

// Hijack lets the caller take over the connection.
// After a call to Hijack(), the HTTP server library
// will not do anything else with the connection.
// It becomes the caller's responsibility to manage
// and close the connection.
func (r Res) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	r.c.pri.cut = true

	hj, ok := r.rw.(http.Hijacker)
	if ok {
		return hj.Hijack()
	}

	return nil, nil, ErrorStr("Connection is not Hijackable")
}

// Flush sends any buffered data to the client.
func (r Res) Flush() {
	fl, ok := r.rw.(http.Flusher)
	if ok {
		fl.Flush()
	}
}

type private struct {
	path       string
	pathAlt    string
	curpath    string
	reswrite   io.Writer
	cut        bool
	firstWrite bool
	session    *SessionAdv
	secure     bool
}

// Strictly Public Variable
type Public struct {
	// Error Code
	Status int
	// Data, useful for storing login credentail
	Data map[string]interface{}
	// Well same as data, but for string data type only! Useful for storing user country code!
	DataStr map[string]string
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
	// ISO-639 and ISO-3166, e.g en-GB (English Great Britain)
	LangCode string
}

// The Framework Structure and Context of User and the Application
type Context struct {
	App *App
	// Request
	Req *http.Request
	// Responce
	Res Res
	// Public Variabless
	Pub Public
	pri private
}

// true if output was sent to client, otherwise false!
func (c *Context) Terminated() bool {
	return c.pri.cut
}

// Signal framework to end user request
func (c *Context) Terminate() {
	c.pri.cut = true
}

func (c *Context) debuginfo() {
	ErrPrintf("%s, %s, %d, %s, %s, ?%s IP:%s, %v",
		c.Req.Proto, c.Req.Method, c.Pub.Status,
		c.Req.Host, c.Http().Path(),
		c.Req.URL.RawQuery, c.Req.RemoteAddr, time.Now())
	ErrPrintln()
}
