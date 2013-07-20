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

type Time struct {
	c *Core
}

func (c *Core) Time() Time {
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

type TimeMiddleware struct {
	Middleware
}

func (t *TimeMiddleware) Html() {
	c := t.C
	// Convert to Default Timezone.
	c.Pub.HtmlFunc["time"] = func(clock time.Time) time.Time {
		return clock.In(c.Pub.TimeLoc)
	}

	// Convert to Timezone
	c.Pub.HtmlFunc["timeZone"] = func(zone string, clock time.Time) time.Time {
		loc, err := time.LoadLocation(zone)
		c.Check(err)
		return clock.In(loc)
	}

	// Format time, leave empty for default
	c.Pub.HtmlFunc["timeFormat"] = func(format string, clock time.Time) string {
		if format == "" {
			format = c.Pub.TimeFormat
		}
		return clock.Format(format)
	}
}
