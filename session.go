package core

import (
	"encoding/gob"
	"os"
	"sync"
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
func (ses *session) hit() {
	ses.Expire = time.Now().Add(SessionExpire)
}

// Session interface
//
// Note: Useful for checking the existant of the session!
type sessionInterface interface {
	getData() interface{}
	getExpire() time.Time
	hit()
}

var sessionMap = struct {
	sync.Mutex
	m map[string]sessionInterface
}{m: map[string]sessionInterface{}}

type SessionHandler interface {
	Set(*Core, interface{})
	Init(*Core)
	Destroy(*Core)
}

// Store Session to Memory
type SessionMemory struct{}

func (_ SessionMemory) Set(c *Core, data interface{}) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	if !sessionExpiryCheckActive {
		sessionExpiryCheckActive = true
		go sessionExpiryCheck()
	}

	sessionCookieName := SessionCookieName.String()

	sesCookie, err := c.Cookie(sessionCookieName).Get()

	if err != nil {
		sesCookie, _ = c.Cookie(sessionCookieName).Value(KeyGen()).SaveRes().Get()
	}

	sessionMap.m[sesCookie.Value] = &session{data, time.Now().Add(SessionExpire)}
}

func deleteSessionFromMap(key string) {
	delete(sessionMap.m, key)
}

func (_ SessionMemory) Init(c *Core) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	sesCookie, err := c.Cookie(SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	if t, ok := sessionMap.m[sesCookie.Value].(*session); ok {
		if time.Now().Unix() < t.getExpire().Unix() {
			c.Pub.Session = t.getData()
			t.hit()
			return
		}
	}

	deleteSessionFromMap(sesCookie.Value)
	c.Cookie(sesCookie.Name).Delete()
}

func (_ SessionMemory) Destroy(c *Core) {
	sessionMap.Lock()
	defer sessionMap.Unlock()

	sesCookie, err := c.Cookie(SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	if _, ok := sessionMap.m[sesCookie.Value].(*session); ok {
		deleteSessionFromMap(sesCookie.Value)
	}

	c.Cookie(sesCookie.Name).Delete()
}

const sessionFileExt = ".wbs"

// Store Session to File.
type SessionFile struct {
	Path string
}

func (se SessionFile) Set(c *Core, data interface{}) {
	sessionCookieName := SessionCookieName.String()
	sesCookie, err := c.Cookie(sessionCookieName).Get()

	if err != nil {
		sesCookie, _ = c.Cookie(sessionCookieName).Value(KeyGen()).SaveRes().Get()
	}

	file, err := os.Create(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Check(err)

	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(&session{data, time.Now().Add(SessionExpire)})
	if err != nil {
		panic(err)
	}
}

func (se SessionFile) Init(c *Core) {
	sesCookie, err := c.Cookie(SessionCookieName.String()).Get()
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
		ses.hit()
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Cookie(sesCookie.Name).Delete()
}

func (se SessionFile) Destroy(c *Core) {
	sesCookie, err := c.Cookie(SessionCookieName.String()).Get()
	if err != nil {
		return
	}

	os.Remove(se.Path + "/" + sesCookie.Value + sessionFileExt)
	c.Cookie(sesCookie.Name).Delete()
}

// Init Session
func (c *Core) initSession() {
	DefaultSessionHandler.Init(c)
}

type Session struct {
	c *Core
}

func (c *Core) Session() Session {
	return Session{c}
}

// Get Session
func (s Session) Get() interface{} {
	return s.c.Pub.Session
}

// Set Session
func (s Session) Set(data interface{}) {
	DefaultSessionHandler.Set(s.c, data)
}

// Destroy Session
func (s Session) Destroy() {
	DefaultSessionHandler.Destroy(s.c)
}

//	Session Expiry Check in a loop
func sessionExpiryCheck() {
	for {
		time.Sleep(SessionExpiryCheckInterval)
		curtime := time.Now()

		sessionMap.Lock()

		if len(sessionMap.m) <= 0 {
			sessionExpiryCheckActive = false
			sessionMap.Unlock()
			break
		}
		for key, value := range sessionMap.m {
			if curtime.Unix() > value.getExpire().Unix() {
				delete(sessionMap.m, key)
			}
		}

		sessionMap.Unlock()
	}
}
