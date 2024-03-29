package core

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"time"
)

func printPanic(buf io.Writer, c *Context, r interface{}, stack []byte) {
	printF := func(format string, a ...interface{}) {
		fmt.Fprintf(buf, format, a...)
	}

	printLn := func(a ...interface{}) {
		fmt.Fprintln(buf, a...)
	}

	printF("\r\n%s, %s, %s, %s, ?%s IP:%s\r\n",
		c.Req.Proto, c.Req.Method,
		c.Req.Host, c.Req.URL.Path,
		c.Req.URL.RawQuery, c.Req.RemoteAddr)

	printF("\r\n%s\r\n\r\n%s", r, stack)

	printLn("\r\nRequest Header:")
	printLn(c.Req.Header)

	c.Req.ParseMultipartForm(c.App.FormMemoryLimit)

	printLn("\r\nForm Values:")
	printLn(c.Req.Form)

	printLn("\r\nForm Values (Multipart):")
	printLn(c.Req.MultipartForm)

	printLn("\r\nTime:")
	printLn(time.Now())
}

// Check for Error
func (c *Context) Check(err error) {
	Check(err)
}

// Check for Error
func Check(err error) {
	if err != nil {
		panic(err)
	}
}

type PanicHandler interface {
	Panic(*Context, interface{}, []byte)
}

// Write error to Stderr.
type PanicConsole struct{}

func (_ PanicConsole) Panic(c *Context, r interface{}, stack []byte) {
	ErrPrint(r, "\r\n", string(stack))
}

const panicFileExt = ".txt"

// Write error to new file.
type PanicFile struct {
	Path string
}

func (p PanicFile) Panic(c *Context, r interface{}, stack []byte) {
	filename := p.Path + fmt.Sprintf("/%d_%d", time.Now().Unix(), time.Now().UnixNano()) + panicFileExt
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer file.Close()
	printPanic(file, c, r, stack)
}

var DefaultPanicHandler PanicHandler = PanicConsole{}

type Errors struct {
	E403 func(c *Context)
	E404 func(c *Context)
	E405 func(c *Context)
	E500 func(c *Context)
}

// Execute Error 403 (Forbidden)
func (c *Context) Error403() {
	c.Pub.Status = 403
	c.Pub.Errors.E403(c)
	c.Terminate()
}

// Execute Error 404 (Not Found)
func (c *Context) Error404() {
	c.Pub.Status = 404
	c.Pub.Errors.E404(c)
	c.Terminate()
}

// Execute Error 405 (Method Not Allowed)
func (c *Context) Error405() {
	c.Pub.Status = 405
	c.Pub.Errors.E405(c)
	c.Terminate()
}

// Execute Error 500 (Internal Server Error)
func (c *Context) Error500() {
	c.Pub.Status = 500
	c.Pub.Errors.E500(c)
	c.Terminate()
}

// Custom String Data Type, Implement error interface.
type ErrorStr string

func (e ErrorStr) Error() string {
	return "Error: " + string(e)
}

// Print formats using the default formats for its operands and writes to standard error output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func ErrPrint(a ...interface{}) {
	fmt.Fprint(os.Stderr, a...)
}

// Printf formats according to a format specifier and writes to standard error output.
// It returns the number of bytes written and any write error encountered.
func ErrPrintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

// Println formats using the default formats for its operands and writes to standard error output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func ErrPrintln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func (c *Context) recover() {
	if r := recover(); r != nil {
		stack := debug.Stack()
		DefaultPanicHandler.Panic(c, r, stack)
		if c.App.Debug {
			c.Pub.Status = 500
			c.Fmt().Println("500 Internal Server Error")
			printPanic(c.Res, c, r, stack)
			return
		}
		c.Error500()
	}
}
