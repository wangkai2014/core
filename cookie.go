package core

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// Chainable version of 'net/http.Cookie'
type Cookie struct {
	core *Core
	c    *http.Cookie
}

// New Cookie
func NewCookie(c *Core, name string) Cookie {
	return Cookie{
		core: c,
		c:    &http.Cookie{Name: name},
	}
}

// Cookie
func (c *Core) Cookie(name string) Cookie {
	return NewCookie(c, name)
}

// Set Value
func (c Cookie) Value(value string) Cookie {
	c.c.Value = value
	return c
}

// Set Path
func (c Cookie) Path(path string) Cookie {
	c.c.Path = path
	return c
}

// Set Domain
func (c Cookie) Domain(domain string) Cookie {
	c.c.Domain = domain
	return c
}

// Set Expiry Time of Cookie.
func (c Cookie) Expires(expires time.Time) Cookie {
	c.c.Expires = expires
	return c
}

// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
func (c Cookie) MaxAge(maxage int) Cookie {
	c.c.MaxAge = maxage
	return c
}

// Make Cookie Secure
func (c Cookie) Secure() Cookie {
	c.c.Secure = true
	return c
}

// Make Cookie Http Only
func (c Cookie) HttpOnly() Cookie {
	c.c.HttpOnly = true
	return c
}

// Get *http.Cookie, if Value is not set it will try to get the Cookie from the User Request!
func (c Cookie) Get() (*http.Cookie, error) {
	if c.c.Value != "" {
		return c.c, nil
	}
	return c.core.Req.Cookie(c.c.Name)
}

// Delete Cookie
func (c Cookie) Delete() Cookie {
	return c.Value("Delete-Me").MaxAge(-1).SaveRes()
}

// Save (Set) Cookie to Response
func (c Cookie) SaveRes() Cookie {
	http.SetCookie(c.core, c.pre(c.c))
	return c
}

// Prepare Cookie
func (c Cookie) pre(cookie *http.Cookie) *http.Cookie {
	if cookie.Path == "" {
		cookie.Path = "/"
	}

	if cookie.Domain != "" {
		return cookie
	}

	cookie.Domain = c.core.Req.Host

	// Split port from address.
	if host, _, err := net.SplitHostPort(cookie.Domain); err == nil {
		cookie.Domain = host
	}

	// Determine if IP address!
	if net.ParseIP(strings.Trim(cookie.Domain, "[]")) != nil {
		cookie.Domain = ""
		return cookie
	}

	// Make sure it's actually a domain name, a domain name has at least one period (.).
	if strings.Count(cookie.Domain, ".") <= 0 {
		cookie.Domain = ""
	}

	return cookie
}

// Save (Add) Cookie to Request, It won't send anything out to the client.
// But it is a useful feature for CSRF protection for example!.
func (c Cookie) SaveReq() Cookie {
	c.core.Req.AddCookie(c.c)
	return c
}

// Set to Expire after an hour
func (c Cookie) Hour() Cookie {
	return c.Expires(time.Now().Add(1 * time.Hour))
}

// Set to Expire after 6 Hourss
func (c Cookie) SixHours() Cookie {
	return c.Expires(time.Now().Add(6 * time.Hour))
}

// Set to Expire after 12 Hourss
func (c Cookie) TwelveHours() Cookie {
	return c.Expires(time.Now().Add(12 * time.Hour))
}

// Set to Expire after 1 Day
func (c Cookie) Day() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 1))
}

// Set to Expire after 1 Weeks
func (c Cookie) Week() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 1*7))
}

// Set to Expire after 2 Weeks
func (c Cookie) TwoWeeks() Cookie {
	return c.Expires(time.Now().AddDate(0, 0, 2*7))
}

// Set to Expire after 1 Month
func (c Cookie) Month() Cookie {
	return c.Expires(time.Now().AddDate(0, 1, 0))
}

// Set to Expire after 3 Months
func (c Cookie) ThreeMonths() Cookie {
	return c.Expires(time.Now().AddDate(0, 3, 0))
}

// Set to Expire after 6 Months
func (c Cookie) SixMonths() Cookie {
	return c.Expires(time.Now().AddDate(0, 6, 0))
}

// Set to Expire after 9 Months
func (c Cookie) NineMonths() Cookie {
	return c.Expires(time.Now().AddDate(0, 9, 0))
}

// Set to Expire after 1 Year
func (c Cookie) Year() Cookie {
	return c.Expires(time.Now().AddDate(1, 0, 0))
}

// Set to Expire after 2 Years
func (c Cookie) TwoYears() Cookie {
	return c.Expires(time.Now().AddDate(2, 0, 0))
}

// Set to Expire after 3 Years
func (c Cookie) ThreeYears() Cookie {
	return c.Expires(time.Now().AddDate(3, 0, 0))
}

// Set to Expire after 4 Years
func (c Cookie) FourYears() Cookie {
	return c.Expires(time.Now().AddDate(4, 0, 0))
}

// Set to Expire after 5 Years
func (c Cookie) FiveYears() Cookie {
	return c.Expires(time.Now().AddDate(5, 0, 0))
}
