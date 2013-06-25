package core

import (
	"io/ioutil"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Convert String to int64
func toInt(number string) (int64, error) {
	return strconv.ParseInt(number, 10, 64)
}

// Convert String to Uint64
func toUint(number string) (uint64, error) {
	return strconv.ParseUint(number, 10, 64)
}

// Convert String to float64
func toFloat(number string) (float64, error) {
	return strconv.ParseFloat(number, 64)
}

func (c *Core) initWriter() {
	if c.Req.Method == "HEAD" {
		c.pri.reswrite = ioutil.Discard
		c.Header().Set("Connection", "close")
		return
	}
	c.pri.reswrite = c.rw
	c.Header().Set("Content-Encoding", "plain")
}

func (c *Core) initTrueHost() {
	switch {
	case c.Req.Header.Get("Host") != "":
		c.Req.Host = c.Req.Header.Get("Host")
	case c.Req.Header.Get("X-Forwarded-Host") != "":
		c.Req.Host = c.Req.Header.Get("X-Forwarded-Host")
	case c.Req.Header.Get("X-Forwarded-Server") != "":
		c.Req.Host = c.Req.Header.Get("X-Forwarded-Server")
	}
	c.Req.URL.Host = c.Req.Host
}

func (c *Core) initTrueRemoteAddr() {
	address := ""

	switch {
	case c.Req.Header.Get("X-Real-Ip") != "":
		address = c.Req.Header.Get("X-Real-Ip")
	case c.Req.Header.Get("X-Forwarded-For") != "":
		address = c.Req.Header.Get("X-Forwarded-For")
	}

	if address != "" {
		c.Req.RemoteAddr = net.JoinHostPort(strings.Trim(address, "[]"), "1234")
	}
}

func (c *Core) initTruePath() {
	switch {
	case c.Req.Header.Get("X-Original-Url") != "":
		// For compatibility with IIS
		urls := strings.Split(c.Req.Header.Get("X-Original-Url"), "?")
		c.Req.URL.Path = urls[0]
		c.pri.path = c.Req.URL.Path

		if len(urls) < 2 {
			return
		}

		c.Req.URL.RawQuery = strings.Join(urls[1:], "?")
	}
}

func (c *Core) initSecure() {
	if c.Req.Header.Get("X-Secure-Mode") != "" {
		c.Req.Proto = "S" + c.Req.Proto
		c.Req.Header.Del("X-Secure-Mode")
	}
}

// Get Remote Address (IP Address) without port number!
func (c *Core) RemoteAddr() string {
	ip, _, _ := net.SplitHostPort(c.Req.RemoteAddr)
	return ip
}

type AtomicString struct {
	sync.Mutex
	s string
}

func NewAtomicString(s string) *AtomicString {
	return &AtomicString{s: s}
}

func (str *AtomicString) String() string {
	str.Lock()
	defer str.Unlock()
	return str.s
}

func (str *AtomicString) Set(s string) {
	str.Lock()
	defer str.Unlock()
	str.s = s
}
