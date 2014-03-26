package core

import (
	"fmt"
	"io"
)

type Fmt struct {
	c *Context
}

func (c *Context) Fmt() Fmt {
	return Fmt{c}
}

// Print formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Print(a ...interface{}) (int, error) {
	return fmt.Fprint(f.c.Res, a...)
}

// Printf formats according to a format specifier and writes to client (http web server or browser).
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Printf(format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f.c.Res, format, a...)
}

// Println formats using the default formats for its operands and writes to client (http web server or browser).
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Println(a ...interface{}) (int, error) {
	return fmt.Fprintln(f.c.Res, a...)
}

// Sprint formats using the default formats for its operands and returns the resulting string.
// Spaces are added between operands when neither is a string.
func (f Fmt) Sprint(a ...interface{}) string {
	return fmt.Sprint(a...)
}

// Sprintf formats according to a format specifier and returns the resulting string.
func (f Fmt) Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

// Sprintln formats using the default formats for its operands and returns the resulting string.
// Spaces are always added between operands and a newline is appended.
func (f Fmt) Sprintln(a ...interface{}) string {
	return fmt.Sprintln(a...)
}

// Fprint formats using the default formats for its operands and writes to w.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Fprint(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprint(w, a...)
}

// Fprintf formats according to a format specifier and writes to w.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Fprintf(w io.Writer, format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(w, format, a...)
}

// Fprintln formats using the default formats for its operands and writes to w.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Fprintln(w io.Writer, a ...interface{}) (int, error) {
	return fmt.Fprintln(w, a...)
}

// Wprint formats using the default formats for its operands and writes to c.Pub.Writers[writerName].
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Wprint(writerName string, a ...interface{}) (int, error) {
	return fmt.Fprint(f.c.IO().W(writerName), a...)
}

// Wprintf formats according to a format specifier and writes to c.Pub.Writers[writerName].
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Wprintf(writerName, format string, a ...interface{}) (int, error) {
	return fmt.Fprintf(f.c.IO().W(writerName), format, a...)
}

// Wprintln formats using the default formats for its operands and writes to c.Pub.Writers[writerName].
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (f Fmt) Wprintln(writerName string, a ...interface{}) (int, error) {
	return fmt.Fprintln(f.c.IO().W(writerName), a...)
}
