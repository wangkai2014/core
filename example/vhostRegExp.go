package main

import (
	"github.com/gorail/core"
)

func Index(c *core.Core) {
	c.Fmt().Print("This is ", c.Pub.Group["Subdomain"], "example.com")
}

func init() {
	vhostRegExp := core.NewVHostRegExp(core.Map{
		`^(?P<Subdomain>[a-zA-Z0-9]*)\.?example\.com`: core.NewRouter().RegisterFunc("^/", Index),
	})

	// Override Main View
	core.MainView = core.RouteHandlerFunc(func(c *core.Core) {
		appMiddlewares := core.AppMiddlewares.Init(c)
		defer appMiddlewares.Post()
		appMiddlewares.Pre()
		if c.CutOut() {
			return
		}

		vhostRegExp.View(c)
	})
}

func main() {
	core.UseMaxCPU()
	core.Check(core.StartHttp(":8080"))
}
