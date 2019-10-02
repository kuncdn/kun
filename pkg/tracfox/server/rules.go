package server

import (
	"net/http"
	"regexp"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"tracfox.io/tracfox/internal/responsewriter"
	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/proxy"
)

// DefaultMaxRetryTimes .
const DefaultMaxRetryTimes = 10

// Rule .
type Rule struct {
	name          string
	methods       []string
	rewiteURIPath string
	reg           *regexp.Regexp
	chain         alice.Chain
	balancer      proxy.Balancer
}

// Match .
func (r *Rule) Match(path string, method string) bool {
	return r.reg.MatchString(path) && util.Contains(r.methods, method)
}

func (r *Rule) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resp := responsewriter.NewBuffer(rw)
	defer resp.FlushAll()
	if len(r.rewiteURIPath) != 0 {
		match := r.reg.FindSubmatchIndex([]byte(req.URL.Path))
		req.URL.Path = string(r.reg.Expand(nil, []byte(r.rewiteURIPath), []byte(req.URL.Path), match))
	}
	for i := 0; i < DefaultMaxRetryTimes; i++ {
		next, err := r.balancer.Elect(req)
		if err != nil {
			glog.Errorln(err)
			http.Error(resp, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
			return
		}
		r.chain.Then(next).ServeHTTP(resp, req)
		if util.IsSuccessCode(resp.Code) {
			return
		}
		glog.Errorf("read the wrong http status code %d, try elect another backend server.", resp.Code)
		resp.Reset()
	}
	glog.Errorln("the maximum number of retries has been exceeded.")
	http.Error(resp, http.StatusText(http.StatusBadGateway), http.StatusBadGateway)
}

// RuleMgr .
type RuleMgr struct {
	rules []*Rule // 扫描等等
}

func (r *RuleMgr) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, rule := range r.rules {
		if rule.Match(req.URL.Path, req.Method) {
			rule.ServeHTTP(rw, req)
			return
		}
	}
	http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
}
