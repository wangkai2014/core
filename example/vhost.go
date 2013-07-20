package main

import (
	"github.com/gorail/core"
)

var app = core.NewApp()

func exampleCom(c *core.Core) {
	c.Fmt().Println("This is example.com")
}

func wwwExampleCom(c *core.Core) {
	c.Fmt().Println("This is dev.example.com")
}

func init() {
	app.DefaultRouter = app.VHost("main").Register(core.Map{
		"example.com":     app.Router("exampleA").RegisterFunc("^/", exampleCom),
		"dev.example.com": app.Router("exampleB").RegisterFunc("^/", wwwExampleCom),
	})
}

func main() {
	core.Check(app.Listen(":8080"))
}
