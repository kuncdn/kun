package proxy

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/config"
)

// NewBalancedReverseProxy .
func NewBalancedReverseProxy(def config.Server) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: createDirector(def),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(def.TCPTimeout) * time.Second,
				KeepAlive: time.Duration(def.TCPKeepAlive) * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConnsPerHost:   def.MaxIdleConnsPerHost,
			MaxIdleConns:          100,
			IdleConnTimeout:       time.Duration(def.IdleConnTimeout) * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
}

func createDirector(def config.Server) func(*http.Request) {
	target, err := url.Parse(def.Target)
	if err != nil {
		panic(err)
	}
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = util.SingleJoiningSlash(target.Path, req.URL.Path)
		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			req.Header.Set("User-Agent", "")
		}
	}
}
