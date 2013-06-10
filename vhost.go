package core

import (
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
	hosts []*vHost
}

// Use host name as string (e.g example.com)
type VHostMap map[string]RouteHandler

// Construct new VHost and Register Map to VHost
func NewVHost(hosts VHostMap) *VHost {
	v := &VHost{}

	v.register(hosts)

	return v
}

func (v *VHost) register(hosts VHostMap) {
	v.Lock()
	defer v.Unlock()
	for host, routerHandler := range hosts {
		v.hosts = append(v.hosts, &vHost{host, routerHandler})
	}
}

func (v *VHost) Add(hosts VHostMap) *VHost {
	v.register(hosts)
	return v
}

func (v *VHost) getHosts() []*vHost {
	v.RLock()
	defer v.RUnlock()
	hosts := []*vHost{}
	hosts = append(hosts, v.hosts...)
	return hosts
}

func (v *VHost) View(c *Core) {
	for _, host := range v.getHosts() {
		if len(host.name) > len(c.Req.Host) {
			continue
		}
		if strings.ToLower(host.name) == strings.ToLower(c.Req.Host[:len(host.name)]) {
			c.RouteDealer(host.route)
			return
		}
	}

	c.Error404()
}

// Use host name regexp as string (e.g. (?P<subdomain>[a-z0-9-_]+)\.example\.com)
type VHostRegExpMap map[string]RouteHandler

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

// Construct VHostRegExp and Register map to VHostRegExp
func NewVHostRegExp(hostmap VHostRegExpMap) *VHostRegExp {
	vh := &VHostRegExp{}
	vh.registerMap(hostmap)
	return vh
}

func (vh *VHostRegExp) getHosts() vHostRegs {
	vh.RLock()
	defer vh.RUnlock()
	hosts := vHostRegs{}
	hosts = append(hosts, vh.vhost...)
	return hosts
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

func (vh *VHostRegExp) registerMap(hostmap VHostRegExpMap) {
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

func (vh *VHostRegExp) Add(hostmap VHostRegExpMap) *VHostRegExp {
	vh.registerMap(hostmap)
	return vh
}

func (vh *VHostRegExp) View(c *Core) {
	for _, host := range vh.getHosts() {
		if !host.RegExpComplied.MatchString(c.Req.Host) {
			continue
		}

		c.pathDealer(host.RegExpComplied, vHostStr(c.Req.Host))

		c.RouteDealer(host.Route)
		return
	}

	c.Error404()
}
