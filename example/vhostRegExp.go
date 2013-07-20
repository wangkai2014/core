package main

import (
	"github.com/gorail/core"
)

var app = core.NewApp()

func Index(c *core.Core) {
	c.Fmt().Print("This is ", c.Pub.Group["Subdomain"], "example.com")
}

func init() {
	app.DefaultRouter = app.VHostRegExp("main").Register(core.Map{
		`^(?P<Subdomain>[a-zA-Z0-9]*)\.?example\.com`: core.NewRouter().RegisterFunc("^/", Index),
	})
}

func main() {
	core.Check(app.Listen(":8080"))
}
