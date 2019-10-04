package server

import (
	"net/http"
	"regexp"
)

// VirtualHost .
type VirtualHost struct {
	virtualHostRegs []*regexp.Regexp
	next            http.Handler
}

// Match .
func (d *VirtualHost) Match(virtualHost string) bool {
	for _, reg := range d.virtualHostRegs {
		if reg.MatchString(virtualHost) {
			return true
		}
	}
	return false
}

func (d *VirtualHost) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	d.next.ServeHTTP(rw, req)
}

// VirtualHostMgr .
type VirtualHostMgr struct {
	virtualHosts []*VirtualHost
}

func (d *VirtualHostMgr) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, virtualHost := range d.virtualHosts {
		if virtualHost.Match(req.Host) {
			virtualHost.ServeHTTP(rw, req)
			return
		}
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
