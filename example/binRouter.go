package main

import (
	"github.com/gorail/core"
)

func Index(c *core.Core) {
	c.Fmt().Print("Index")
}

func Example(c *core.Core) {
	c.Fmt().Print("Example")
}

func SubExample(c *core.Core) {
	c.Fmt().Print("SubExample")
}

func init() {
	core.BinRoute.RootDirFunc(Index).RegisterMap(core.Map{
		"example": core.NewBinRouter().RootDirFunc(Example).RegisterFuncMap(core.FuncMap{
			"subexample": SubExample,
		}),
	})
}

func main() {
	core.SetBinRouteToMainView()
	core.Check(core.StartHttp(":8080"))
}
