package core

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleCore_Json(c *Core) {
	// Prepare structure
	data := struct {
		Title string `json:"title"`
		Life  int    `json:"life"`
	}{}

	// Decode from Request Body
	c.Json().DecodeReqBody(&data)

	// Send it back to the client
	c.Json().Send(data)
}

func TestJson(t *testing.T) {
	type jsonTest struct {
		Title string `json:"title"`
	}

	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(c *Core) {
		jsonTest2 := jsonTest{}

		c.Json().DecodeReqBody(&jsonTest2)

		if jsonTest2.Title != "Hello World!" {
			t.Fail()
		}

		jsonTest2.Title = "Hola!"

		c.Json().Send(jsonTest2)
	})

	jsonTest1 := jsonTest{"Hello World!"}

	b, err := json.Marshal(jsonTest1)
	Check(err)

	ts := httptest.NewServer(App)
	defer ts.Close()

	req, err := http.NewRequest("POST", ts.URL, bytes.NewReader(b))
	Check(err)

	res, err := http.DefaultClient.Do(req)
	Check(err)

	err = json.NewDecoder(res.Body).Decode(&jsonTest1)
	Check(err)

	if jsonTest1.Title != "Hola!" {
		t.Fail()
	}
}
