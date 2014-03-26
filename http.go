package core

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"net/url"
)

type Http struct {
	c *Context
}

func (c *Context) Http() Http {
	return Http{c}
}

// Set Cookie
func (h Http) SetCookie(cookie *http.Cookie) {
	http.SetCookie(h.c.Res, cookie)
}

// Get Cookie
func (h Http) GetCookie(name string) (*http.Cookie, error) {
	return h.c.Req.Cookie(name)
}

// Execute Handler
func (h Http) Exec(handler http.Handler) {
	http.StripPrefix(h.c.pri.curpath, handler).ServeHTTP(h.c.Res, h.c.Req)
}

// Execute Function
func (h Http) ExecFunc(handler http.HandlerFunc) {
	h.Exec(handler)
}

// ServeFile replies to the request with the contents of the named file or directory.
func (h Http) ServeFile(name string) {
	http.ServeFile(h.c.Res, h.c.Req, name)
}

// Get issues a GET to the specified URL.  If the response is one of the following
// redirect codes, Get follows the redirect, up to a maximum of 10 redirects:
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// An error is returned if there were too many redirects or if there
// was an HTTP protocol error. A non-2xx response doesn't cause an
// error.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// Get is a wrapper around http.DefaultClient.Get.
func (h Http) Get(url string) (*http.Response, error) {
	return http.Get(url)
}

// Get Body from Url as Bytes
func (h Http) GetBytes(url string) ([]byte, error) {
	resp, err := h.Get(url)
	if err != nil {
		return nil, err
	}
	b := &bytes.Buffer{}
	defer b.Reset()
	io.Copy(b, resp.Body)
	return b.Bytes(), nil
}

// Head issues a HEAD to the specified URL.  If the response is one of the
// following redirect codes, Head follows the redirect after calling the
// Client's CheckRedirect function.
//
//    301 (Moved Permanently)
//    302 (Found)
//    303 (See Other)
//    307 (Temporary Redirect)
//
// Head is a wrapper around http.DefaultClient.Head
func (h Http) Head(url string) (*http.Response, error) {
	return http.Head(url)
}

// Post issues a POST to the specified URL.
//
// Caller should close resp.Body when done reading from it.
//
// Post is a wrapper around http.DefaultClient.Post
func (h Http) Post(url string, bodyType string, body io.Reader) (*http.Response, error) {
	return http.Post(url, bodyType, body)
}

// PostForm issues a POST to the specified URL, with data's keys and
// values URL-encoded as the request body.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// PostForm is a wrapper around http.DefaultClient.PostForm
func (h Http) PostForm(url string, data url.Values) (*http.Response, error) {
	return http.PostForm(url, data)
}

// ReadResponse reads and returns an HTTP response from r.  The
// req parameter specifies the Request that corresponds to
// this Response.  Clients must call resp.Body.Close when finished
// reading resp.Body.  After that call, clients can inspect
// resp.Trailer to find key/value pairs included in the response
// trailer.
func (h Http) ReadResponse(r *bufio.Reader, req *http.Request) (*http.Response, error) {
	return http.ReadResponse(r, req)
}

// Path
func (h Http) Path() string {
	return h.c.pri.curpath + h.c.pri.path
}
