package core

import (
	"bytes"
	"io"
)

// Buffer Shortcut! (c.Pub.Readers and c.Pub.Writers)
type IO struct {
	c *Context
}

// Buffer Shortcut!
func (c *Context) IO() IO {
	return IO{c}
}

// Get Writer
func (io_ IO) W(name string) io.Writer {
	return io_.c.Pub.Writers[name]
}

// Get Reader
func (io_ IO) R(name string) io.Reader {
	return io_.c.Pub.Readers[name]
}

// Push from Reader to Writer
func (io_ IO) PushRtoW(readerName, writerName string) {
	r := io_.R(readerName)
	w := io_.W(writerName)
	if r == nil || w == nil {
		return
	}
	io.Copy(w, r)
}

// Push Content to Writer
func (io_ IO) Push(writerName string, content []byte) {
	w := io_.W(writerName)
	if w == nil {
		return
	}
	io.Copy(w, bytes.NewReader(content))
}

// Push Content to Writer as string
func (io_ IO) PushStr(writerName, content string) {
	io_.Push(writerName, []byte(content))
}

// Push Content to Writer from io.Reader
func (io_ IO) PushReader(writerName string, r io.Reader) {
	w := io_.W(writerName)
	if w == nil {
		return
	}
	io.Copy(w, r)
}

// Push Content direct to Client
func (io_ IO) PushToClient(content []byte) {
	io.Copy(io_.c.Res, bytes.NewReader(content))
}

// Push Content direct to Client as string
func (io_ IO) PushToClientStr(content string) {
	io_.PushToClient([]byte(content))
}

// Push Content direct to Client from io.Reader
func (io_ IO) PushToClientReader(r io.Reader) {
	io.Copy(io_.c.Res, r)
}

// Pull Content from a Reader as []byte
func (io_ IO) Pull(readerName string) []byte {
	r := io_.R(readerName)
	if r == nil {
		return nil
	}
	b := []byte{}
	for {
		bb := make([]byte, 1024)
		num, _ := r.Read(bb)
		if num == 0 {
			break
		}
		b = append(b, bb[:num]...)
	}
	return b
}

// Pull Content from a Reader as String
func (io_ IO) PullStr(readerName string) string {
	if io_.R(readerName) == nil {
		return ""
	}
	return string(io_.Pull(readerName))
}

// Pull Content from A Reader to io.Writer
func (io_ IO) PullWriter(readerName string, w io.Writer) {
	if io_.R(readerName) == nil {
		return
	}
	io.Copy(w, io_.R(readerName))
}

// Push Io
func (io_ IO) PushIO(dst io.Writer, src io.Reader) (int64, error) {
	return io.Copy(dst, src)
}
