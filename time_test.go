package core

func ExampleCore_Time(c *Core) {
	// Get current time
	curtime := c.Time().Now()

	// Output to client
	c.Fmt().Println(curtime)

	// Set Timezone on user request level
	c.Time().SetZone("Europe/London")
}
