/*
Copyright 2019 The Tracfox Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
