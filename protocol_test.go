package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type ProtocolDummy struct {
	Protocol
}

func (pr *ProtocolDummy) Http() {
	pr.C.Pub.Group.Set("protocol", "HTTP")
}

func (pr *ProtocolDummy) Https() {
	pr.C.Pub.Group.Set("protocol", "HTTPS")
}

func TestProtocol(t *testing.T) {
	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Core) {

		protocol := func() string {
			return c.Pub.Group.Get("protocol")
		}

		c.RouteDealer(&ProtocolDummy{})

		if protocol() != "HTTP" {
			t.Fail()
		}

		c.pri.secure = true

		c.RouteDealer(&ProtocolDummy{})

		if protocol() != "HTTPS" {
			t.Fail()
		}

	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
