package core

import (
	"sync"
)

var (
	/*
		Middleware
	*/
	MainMiddlewares = NewMiddlewares()

	_routeAsserter = struct {
		sync.RWMutex
		ro []RouteAsserter
	}{ro: []RouteAsserter{}}
)
