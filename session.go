package core

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"os"
	"time"
)

// Structure of Session
type session struct {
	Data   interface{}
	Expire time.Time
}

// Get the session data!
func (ses *session) getData() interface{} {
	return ses.Data
}

// Returns the time of expiry
func (ses *session) getExpire() time.Time {
	return ses.Expire
}

// Reset Expiry Time to 20 minutes in advanced!
func (ses *session) hit(c *Context) {
	ses.Expire = time.Now().Add(c.App.SessionExpire)
}

type sessionStateless struct {
	Data interface{}
}

func init() {
	gob.Register(&sessionStateless{})
}

// Session interface
//
// Note: Useful for checking the existant of the session!
type sessionInterface interface {
	getData() interface{}
	getExpire() time.Time
	hit(*Context)
}

type SessionHandler interface {
	Set(*Context, interface{})
	Init(*Context)
	Destroy(*Context)
}

// Store Session to Cookie
type SessionStateless struct{}

func (_ SessionStateless) Set(c *Context, data interface{}) {
	sessionCookieName := c.App.SessionCookieName.String()

	s := sessionStateless{
		Data: data,
	}

	buf := &bytes.Buffer{}
	defer buf.Reset()
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s)
	Check(err)

	c.Cookie(sessionCookieName).Value(base64.URLEncoding.EncodeToString(buf.Bytes())).Expires(time.Now().Add(c.App.SessionExpire)).SaveRes()
}

func (_ SessionStateless) Init(c *Context) {
	sessionCookieName := c.App.SessionCookieName.String()

	cookie, err := c.Cookie(sessionCookieName).Get()
	if err != nil {
		return
	}

	s := sessionStateless{}

	buf := &bytes.Buffer{}
	defer buf.Reset()
	b, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		c.Cookie(sessionCookieName).Delete()
		return
	}
	buf.Write(b)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&s)
	if err != nil {
		c.Cookie(sessionCookieName).Delete()
		return
	}

	c.Pub.Session = s.Data
	c.Cookie(sessionCookieName).Value(cookie.Value).Expires(time.Now().Add(c.App.SessionExpire)).SaveRes()
}

func (_ SessionStateless) Destroy(c *Context) {
	sessionCookieName := c.App.SessionCookieName.String()
	c.Cookie(sessionCookieName).Delete()
}

// Store Session to Memory
type SessionMemory struct{}

func (_ SessionMemory) Set(c *Context, data interface{}) {
	c.App.sessionMapSync.Lock()
	defer c.App.sessionMapSync.Unlock()

	if !c.App.sessionExpireCheckActive {
		c.App.sessionExpireCheckActive = true
		go c.App.sessionExpiryCheck()
	}

	sessionCookieName := c.App.SessionCookieName.String()

	sesCookie, err := c.Cookie(sessionCookieName).Get()

	if err != nil {
		sesCookie, _ = c.Cookie(sessionCookieName).Value(KeyGen()).SaveRes().Get()
	}

	c.App.sessionMap[sesCookie.Value] = &session{data, time.Now().Add(c.App.SessionExpire)}
}

func (_ SessionMemory) Init(c *Context) {
	c.App.sessionMapSync.Lock()
	defer c.App.sessionMapSync.Unlock()

	sesCookie, err := c.Cookie(c.App.SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	if t, ok := c.App.sessionMap[sesCookie.Value].(*session); ok {
		if time.Now().Unix() < t.getExpire().Unix() {
			c.Pub.Session = t.getData()
			t.hit(c)
			return
		}
	}

	delete(c.App.sessionMap, sesCookie.Value)
	c.Cookie(sesCookie.Name).Delete()
}

func (_ SessionMemory) Destroy(c *Context) {
	c.App.sessionMapSync.Lock()
	defer c.App.sessionMapSync.Unlock()

	sesCookie, err := c.Cookie(c.App.SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	if _, ok := c.App.sessionMap[sesCookie.Value].(*session); ok {
		delete(c.App.sessionMap, sesCookie.Value)
	}

	c.Cookie(sesCookie.Name).Delete()
}

const sessionFileExt = ".wbs"

// Store Session to File.
type SessionFile struct {
	Path string
}

func (se SessionFile) Set(c *Context, data interface{}) {
	sessionCookieName := c.App.SessionCookieName.String()
	sesCookie, err := c.Cookie(sessionCookieName).Get()

	if err != nil {
		sesCookie, _ = c.Cookie(sessionCookieName).Value(KeyGen()).SaveRes().Get()
	}

	file, err := os.Create(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Check(err)

	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(&session{data, time.Now().Add(c.App.SessionExpire)})
	if err != nil {
		panic(err)
	}
}

func (se SessionFile) Init(c *Context) {
	sesCookie, err := c.Cookie(c.App.SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	file, err := os.Open(se.Path + "/" + sesCookie.Value + sessionFileExt)
	if err != nil {
		return
	}
	defer file.Close()
	dec := gob.NewDecoder(file)

	ses := &session{}

	err = dec.Decode(&ses)
	Check(err)

	if time.Now().Unix() < ses.getExpire().Unix() {
		c.Pub.Session = ses.getData()
		ses.hit(c)
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Cookie(sesCookie.Name).Delete()
}

func (se SessionFile) Destroy(c *Context) {
	sesCookie, err := c.Cookie(c.App.SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Cookie(sesCookie.Name).Delete()
}

// Init Session
func (c *Context) initSession() {
	c.App.SessionHandler.Init(c)
}

type Session struct {
	c *Context
}

func (c *Context) Session() Session {
	return Session{c}
}

// Get Session
func (s Session) Get() interface{} {
	return s.c.Pub.Session
}

// Set Session
func (s Session) Set(data interface{}) {
	s.c.App.SessionHandler.Set(s.c, data)
}

// Destroy Session
func (s Session) Destroy() {
	s.c.App.SessionHandler.Destroy(s.c)
}

//	Session Expiry Check in a loop
func (app *App) sessionExpiryCheck() {
	for {
		time.Sleep(app.SessionExpireCheckInterval)
		curtime := time.Now()

		app.sessionMapSync.Lock()

		if len(app.sessionMap) <= 0 {
			app.sessionExpireCheckActive = false
			app.sessionMapSync.Unlock()
			break
		}
		for key, value := range app.sessionMap {
			if curtime.Unix() > value.getExpire().Unix() {
				delete(app.sessionMap, key)
			}
		}

		app.sessionMapSync.Unlock()
	}
}
