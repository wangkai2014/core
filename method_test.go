package core

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type MethodDummy struct {
	Method
}

func (me *MethodDummy) Prepare() {
	me.C.Pub.Group.Set("prepare", "PREPARE")
}

func (me *MethodDummy) Ws() {
	me.C.Pub.Group.Set("ws", "WS")
}

func (me *MethodDummy) Ajax() {
	me.C.Pub.Group.Set("ajax", "AJAX")
}

func (me *MethodDummy) Finish() {
	me.C.Pub.Group.Set("finish", "FINISH")
}

func (me *MethodDummy) Get() {
	me.C.Pub.Group.Set("method", "GET")
}

func (me *MethodDummy) Post() {
	me.C.Pub.Group.Set("method", "POST")
}

func (me *MethodDummy) Delete() {
	me.C.Pub.Group.Set("method", "DELETE")
}

func (me *MethodDummy) Put() {
	me.C.Pub.Group.Set("method", "PUT")
}

func (me *MethodDummy) Patch() {
	me.C.Pub.Group.Set("method", "PATCH")
}

func (me *MethodDummy) Options() {
	me.C.Pub.Group.Set("method", "OPTIONS")
}

func TestMethod(t *testing.T) {
	App := NewApp()

	App.Debug = true

	App.TestView = RouteHandlerFunc(func(w *Context) {

		prepare := func() string {
			return w.Pub.Group.Get("prepare")
		}

		ws := func() string {
			return w.Pub.Group.Get("ws")
		}

		ajax := func() string {
			return w.Pub.Group.Get("ajax")
		}

		finish := func() string {
			return w.Pub.Group.Get("finish")
		}

		method := func() string {
			return w.Pub.Group.Get("method")
		}

		w.RouteDealer(&MethodDummy{})

		if prepare() != "PREPARE" {
			t.Fail()
		}

		if ws() == "WS" {
			t.Fail()
		}

		if ajax() == "AJAX" {
			t.Fail()
		}

		if finish() != "FINISH" {
			t.Fail()
		}

		if method() != "GET" {
			t.Fail()
		}

		slices := []string{"POST", "DELETE", "PUT", "PATCH", "OPTIONS"}

		for _, value := range slices {
			w.Req.Method = value
			w.RouteDealer(&MethodDummy{})
			if method() != value {
				t.Fail()
			}
		}

	})

	ts := httptest.NewServer(App)
	defer ts.Close()

	http.Get(ts.URL)
}
