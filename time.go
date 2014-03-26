package core

import (
	"time"
)

// Set Timezone on global level
func (app *App) SetTimeZone(zone string) {
	var err error
	app.TimeLoc, err = time.LoadLocation(zone)
	Check(err)
}

// Get Current Time
func CurTime() time.Time {
	return time.Now()
}

// Time
type Time struct {
	c *Context
}

func (c *Context) Time() Time {
	return Time{c}
}

// Set Timezone on user request level
func (t Time) SetZone(zone string) {
	var err error
	t.c.Pub.TimeLoc, err = time.LoadLocation(zone)
	t.c.Check(err)
}

// Get Current Time
func (t Time) Now() time.Time {
	return CurTime()
}
