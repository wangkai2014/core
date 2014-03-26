package main

import (
	"github.com/gorail/core"
)

var app = core.NewApp()

func Index(c *core.Context) {
	c.Fmt().Print("Index")
}

func Example(c *core.Context) {
	c.Fmt().Print("Example")
}

func SubExample(c *core.Context) {
	c.Fmt().Print("SubExample")
}

func init() {
	app.DefaultRouter = app.DirRouter("main").RootDirFunc(Index).RegisterMap(core.Map{
		"example": app.DirRouter("example").RootDirFunc(Example).RegisterFuncMap(core.FuncMap{
			"subexample": SubExample,
		}),
	})
}

func main() {
	core.Check(app.Listen(":8080"))
}
