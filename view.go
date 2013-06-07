package core

// Default Index View
func index(c *Core) {
	c.Fmt().Print("<h1>Hello World!</h1>")
}

// Push Index View to Router
func init() {
	Route.RegisterFunc("^/$", index)
}
