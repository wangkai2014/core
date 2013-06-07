package core

func ExampleCore_Cookie(c *Core) {
	// Setting a Cookie, Expires in a Month.
	c.Cookie("Example").Value("Example").Month().SaveRes()

	// Setting a Cookie, Secure and Http Only.
	c.Cookie("Example").Value("Example").Month().HttpOnly().Secure().SaveRes()

	// Get Cookie from User Request, just omit 'Value' simple as that! Should return *http.Cookie and error
	cookie, err := c.Cookie("Example").Get()

	if err != nil {
		// If the cookie does not exist just set a new cookie and get it.
		cookie, _ = c.Cookie("Example").Value("Example").Month().SaveRes().Get()
	}

	// Delete a Cookie
	c.Cookie(cookie.Name).Delete()

	// Pretty slick, don't you think?
}
