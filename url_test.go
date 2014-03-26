package core

func ExampleContext_Url(c *Context) {
	// Get Absolute Url
	url := c.Url().Absolute("/")

	// Get Absolute Url (Https)
	url = c.Url().AbsoluteHttps("/")

	// Output Url to Client
	c.Fmt().Println(url)
}
