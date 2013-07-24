package core

import (
	"net"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type vHost struct {
	name  string
	route RouteHandler
}

// VHost Data Type, Implement RouteHandler interface.
type VHost struct {
	sync.RWMutex
	hosts map[string]*vHost
}

// Construct New VHost!
func NewVHost() *VHost {
	return &VHost{}
}

func (v *VHost) register(hosts Map) {
	v.Lock()
	defer v.Unlock()
	if v.hosts == nil {
		v.hosts = map[string]*vHost{}
	}
	for host, routerHandler := range hosts {
		v.hosts[host] = &vHost{host, routerHandler}
	}
}

// Use host name as string (e.g example.com)
func (v *VHost) Register(hosts Map) *VHost {
	v.register(hosts)
	return v
}

func (v *VHost) View(c *Core) {
	curHostName := c.Req.Host
	if host, _, err := net.SplitHostPort(curHostName); err == nil {
		curHostName = host
	}
	curHostName = strings.ToLower(curHostName)

	v.RLock()
	host := v.hosts[curHostName]
	v.RUnlock()

	if host == nil {
		c.Error404()
		v.RUnlock()
		return
	}

	c.RouteDealer(host.route)
	return
}

type vHostRegExpItem struct {
	RegExp         string
	RegExpComplied *regexp.Regexp
	Route          RouteHandler
}

type vHostRegs []*vHostRegExpItem

func (vh vHostRegs) Len() int {
	return len(vh)
}

func (vh vHostRegs) Less(i, j int) bool {
	return vh[i].RegExp < vh[j].RegExp
}

func (vh vHostRegs) Swap(i, j int) {
	vh[i], vh[j] = vh[j], vh[i]
}

// VHostRegExp Data type, Implement RouteHandler interface.
type VHostRegExp struct {
	sync.RWMutex
	vhost vHostRegs
}

// Construct VHostRegExp
func NewVHostRegExp() *VHostRegExp {
	return &VHostRegExp{}
}

func (vh *VHostRegExp) register(RegExpRule string, routeHandler RouteHandler) {
	for _, host := range vh.vhost {
		if host.RegExp == RegExpRule {
			host.Route = routeHandler
			return
		}
	}

	vh.vhost = append(vh.vhost, &vHostRegExpItem{RegExpRule, regexp.MustCompile(RegExpRule), routeHandler})
}

func (vh *VHostRegExp) registerMap(hostmap Map) {
	vh.Lock()
	defer vh.Unlock()

	if vh.vhost == nil {
		vh.vhost = vHostRegs{}
	}

	for rule, route := range hostmap {
		vh.register(rule, route)
	}

	sort.Sort(vh.vhost)
}

// Use host name regexp as string (e.g. (?P<subdomain>[a-z0-9-_]+)\.example\.com)
func (vh *VHostRegExp) Register(hostmap Map) *VHostRegExp {
	vh.registerMap(hostmap)
	return vh
}

func (vh *VHostRegExp) View(c *Core) {
	for _, host := range vh.vhost {
		if !host.RegExpComplied.MatchString(c.Req.Host) {
			continue
		}

		c.pathDealer(host.RegExpComplied, vHostStr(c.Req.Host))

		c.RouteDealer(host.Route)
		return
	}

	c.Error404()
}
