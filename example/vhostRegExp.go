package main

import (
	"github.com/gorail/core"
)

func Index(c *core.Core) {
	c.Fmt().Print("This is ", c.Pub.Group["Subdomain"], "example.com")
}

func init() {
	core.VHostsRegExp.Register(core.Map{
		`^(?P<Subdomain>[a-zA-Z0-9]*)\.?example\.com`: core.NewRouter().RegisterFunc("^/", Index),
	})
}

func main() {
	core.UseMaxCPU()
	core.SetVHostsRegExpToMainView()
	core.Check(core.StartHttp(":8080"))
}
