package main

import (
	"github.com/gorail/core"
)

var app = core.NewApp()

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
	app.Router("main").RegisterMap(core.Map{
		`^/$`: &Index{},
		`^/(?P<Id>[0-9-]+)/?$`: &Index{},
	})
}

func main() {
	core.Check(app.Listen(":8080"))
}
