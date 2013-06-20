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
	core.VHosts.Add(core.Map{
		"example.com":     core.NewRouter().RegisterFunc("^/", exampleCom),
		"dev.example.com": core.NewRouter().RegisterFunc("^/", wwwExampleCom),
	})
}

func main() {
	core.UseMaxCPU()
	core.SetVHostsToMainView()
	core.Check(core.StartHttp(":8080"))
}
