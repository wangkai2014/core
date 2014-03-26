package core

func ExampleContext_Header(c *Context) {
	// Set Header
	c.Res.Header().Set("Location", "/world/")
}
