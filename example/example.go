package main

import (
	"github.com/gorail/core"
	"html/template"
)

// Initialise new web application.
var app = core.NewApp()

// Create a RESTful controller with HTML5 buffer system.
type Index struct {
	// There also core.Verb, which exclude the HTML5 Buffer System, just the RESTful controller.
	core.VerbHtml5
	// Id will be populated automatically depending on the url selected, see router below!
	Id int64
}

// Pre Execute before http verb.
func (in *Index) Prepare() {
	if in.Id <= 0 {
		in.Id = 1
	}

	in.RegOnInitFunc(core.HtmlAttrLang("en"), core.HtmlAttrDir("ltr"))

	in.RegOnFinishFunc(func(h core.HtmlPrinter, c *core.Context) {
		// BodyJs comes after Footer, before </body>.
		h.BodyJs(`<script type="text/javascript">
alert("Hello World");
</script>
`)
	})
}

// Execute on Get verb and only on Get verb.
func (in *Index) Get() {
	// Print "Hello World" to <title>...</title>
	in.Title("Hello World")
	// Print Format to Body Content, comes after Header and before Footer.
	in.BodyContentF(`<h1>%s</h1>
<p>%d</p>
`, in.GetBuffer().Title().String(), in.Id)

	// If you want to use a html template system, no problem, just bring your own template system.
	tmpl, _ := template.New("world").Parse(`<h2>{{.}}</h2>
`)
	tmpl.Execute(in.GetBuffer().BodyContent(), in.GetBuffer().Title().String())
}

// Execute on Put verb and only on Put verb
func (in *Index) Put() {
	/*
		The html buffers system will not be initialise unless you call one of the methods below,
		than it will initialise automatically.

		in.GetBuffer
		in.HtmlAttr
		in.HtmlAttrF
		in.HtmlAttrLn
		in.Title
		in.TitleF
		in.TitleLn
		in.Head
		in.HeadF
		in.HeadLn
		in.BodyAttr
		in.BodyAttrF
		in.BodyAttrLn
		in.BodyHeader
		in.BodyHeaderF
		in.BodyHeaderLn
		in.BodyContent
		in.BodyContentF
		in.BodyContentLn
		in.BodyFooter
		in.BodyFooterF
		in.BodyFooterLn
		in.BodyJs
		in.BodyJsF
		in.BodyJsLn

		In some case scenario you just don't need to use the html buffer system.
	*/

	in.C.Fmt().Print("Hello World")
}

func init() {
	// Setup Url Router
	app.Router("main").RegisterMap(core.Map{
		`^/$`: &Index{},
		`^/(?P<Id>[0-9-]+)/?$`: &Index{},
		`^/world`:              app.Router("world"),
	})

	app.Router("world").RegisterMap(core.Map{
		`^/$`: &Index{},
		`^/(?P<Id>[0-9-]+)/?$`: &Index{},
	})

	/*
		If you are dead serious about performance, you can use a thread-safe hash table (map) based router
		'DirRouter' instead of 'Router', 'DirRouter' has a constant time complexity ( O(1) )
		while RegExp based 'Router' has a linear time complexity ( O(n) ).

		app.DefaultRouter = app.DirRouter("main").Root(&Index{}).Asterisk(&Index{}).Group("Id").Register("world", app.DirRouter("world"))

		app.DirRouter("world").Root(&Index{}).Asterisk(&Index{})
	*/

	// Setup Url Reverse, it is seperate from the router system.
	app.URLRev.RegisterMap(core.URLReverseMap{
		"index":            "/",
		"index_page":       "/%d",
		"index_world":      "/world/",
		"index_world_page": "/world/%d",
	})
}

// Start server
func main() {
	core.Check(app.Listen(":8080"))
}
