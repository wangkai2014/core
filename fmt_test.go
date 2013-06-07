package core

func ExampleCore_Fmt(c *Core) {
	// Print to Client
	c.Fmt().Print("Hello", " World")

	// Print to Client with New Line
	c.Fmt().Println("Hello World")

	// Print to Client (Format)
	c.Fmt().Printf("%s %s", "Hello", "World")
}
