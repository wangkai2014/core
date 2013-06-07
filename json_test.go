package core

func ExampleCore_Json(c *Core) {
	// Prepare structure
	data := struct {
		Title string `json:"title"`
		Life  int    `json:"life"`
	}{}

	// Decode from Request Body
	c.Json().DecodeReqBody(&data)

	// Send it back to the client
	c.Json().Send(data)
}
