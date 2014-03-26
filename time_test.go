package core

func ExampleContext_Time(c *Context) {
	// Get current time
	curtime := c.Time().Now()

	// Output to client
	c.Fmt().Println(curtime)

	// Set Timezone on user request level
	c.Time().SetZone("Europe/London")
}
