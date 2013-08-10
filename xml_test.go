package core

import (
	"bytes"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleCore_Xml(c *Core) {
	// Prepare structure
	data := struct {
		Title string `xml:"title,attr"`
		Life  int    `xml:"life"`
	}{}

	// Decode from Request Body
	c.Xml().DecodeReqBody(&data)

	// Send it back to the client
	c.Xml().Send(data)
}

func TestXml(t *testing.T) {
	type xmlTest struct {
		Title string `xml:"title"`
	}

	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Core) {
		xmlTest2 := xmlTest{}

		c.Xml().DecodeReqBody(&xmlTest2)

		if xmlTest2.Title != "Hello World!" {
			t.Fail()
		}

		xmlTest2.Title = "Hola!"

		c.Xml().Send(xmlTest2)
	})

	xmlTest1 := xmlTest{"Hello World!"}

	b, err := xml.Marshal(xmlTest1)
	Check(err)

	ts := httptest.NewServer(App)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(b))
	Check(err)

	res, err := http.DefaultClient.Do(req)
	Check(err)

	err = xml.NewDecoder(res.Body).Decode(&xmlTest1)
	Check(err)

	if xmlTest1.Title != "Hola!" {
		t.Fail()
	}
}
