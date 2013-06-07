package core

import (
	"compress/flate"
	"compress/gzip"
	"strings"
)

// Init Compression Buffer (Call Before Writing to Client)
func (c *Core) InitCompression() {
	if c.Req.Method == "HEAD" {
		return
	}

	if c.Req.Header.Get("Connection") == "Upgrade" {
		return
	}

	for _, encoding := range strings.Split(c.Req.Header.Get("Accept-Encoding"), ",") {
		encoding = strings.TrimSpace(strings.ToLower(encoding))
		switch encoding {
		case "gzip":
			c.pri.reswrite = gzip.NewWriter(c.rw)
			c.Header().Set("Content-Encoding", encoding)
			return
		case "deflate":
			c.pri.reswrite, _ = flate.NewWriter(c.rw, flate.DefaultCompression)
			c.Header().Set("Content-Encoding", encoding)
			return
		}
	}
}

func (c *Core) closeCompression() {
	switch t := c.pri.reswrite.(type) {
	case *gzip.Writer:
		t.Close()
	case *flate.Writer:
		t.Close()
	}
}
