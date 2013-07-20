package core

func ExampleCore_Header(c *Core) {
	// Set Header
	c.Header().Set("Location", "/world/")
}
