package main

import (
	"github.com/gorail/core"
)

type Index struct {
	core.Method
	Id int64
}

func (in *Index) Prepare() {
	if in.Id <= 0 {
		in.Id = 1
	}
}

func (in *Index) Get() {
	const htmlstr = `<h1>Hello World</h1>
	<p>Page: {{.Id}}</p>`
	in.C.Html().RenderSend(htmlstr, in)
}

func init() {
	core.Route.RegisterMap(core.RouteHandlerMap{
		`^/$`: &Index{},
		`^/(?P<Id>[0-9-]+)/?$`: &Index{},
	})
}

func main() {
	core.UseMaxCPU()
	core.Check(core.StartHttp(":8080"))
}