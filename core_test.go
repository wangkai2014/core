package core

import (
	"net"
)

func ExampleCore_Header(c *Core) {
	// Set Header
	c.Header().Set("Location", "/world/")
}

func ExampleStartHttp() {
	// Start Http Server on tcp port 1234, while also checking for errors!
	Check(StartHttp(":1234"))
}

func ExampleStartHttpTLS() {
	// Start Https Server on tcp port 1234, while also checking for errors!
	Check(StartHttpTLS(":1234", "Cert Str", "Key Str"))
}

func ExampleStartFastCGI() {
	// Start FastCGI for Shared Servers, while also checking for errors!
	Check(StartFastCGI(nil))

	// Start FastCGI on tcp port 1234
	l, err := net.Listen("TCP", ":1234")
	Check(err)
	Check(StartFastCGI(l))
}

func ExampleStartCGI() {
	// Start CGI while checking for errors!
	Check(StartCGI())
}
