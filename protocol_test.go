package core

import (
	"net/http"
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
	c := &Core{
		Pub: Public{
			Group: Group{},
		},
		Req: &http.Request{
			Proto:  "HTTP/1.1",
			Header: http.Header{},
		},
		pri: private{
			cut:    false,
			secure: false,
		},
	}

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
}
