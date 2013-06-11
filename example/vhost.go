package main

import (
	"github.com/gorail/core"
)

func exampleCom(c *core.Core) {
	c.Fmt().Println("This is example.com")
}

func wwwExampleCom(c *core.Core) {
	c.Fmt().Println("This is dev.example.com")
}

func init() {
	vhost := core.NewVHost(core.Map{
		"example.com":     core.NewRouter().RegisterFunc("^/", exampleCom),
		"dev.example.com": core.NewRouter().RegisterFunc("^/", wwwExampleCom),
	})

	// Override Main View
	core.MainView = core.RouteHandlerFunc(func(c *core.Core) {
		appMiddlewares := core.AppMiddlewares.Init(c)
		defer appMiddlewares.Post()
		appMiddlewares.Pre()
		if c.CutOut() {
			return
		}

		vhost.View(c)
	})
}

func main() {
	core.UseMaxCPU()
	core.Check(core.StartHttp(":8080"))
}
