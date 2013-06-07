package core

func ExampleCore_Xml(c *Core) {
	// Prepare structure
	data := struct {
		Title string `xml:"title,attr"`
		Life  int    `xml:"life"`
	}{}

	// Decode from Request Body
	c.Xml().DecodeReqBody(&data)

	// Send it back to the client
	c.Xml().Send(data)
}
