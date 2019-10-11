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

package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"tracfox.io/tracfox/internal/responsewriter"
)

type loggerConstructor struct {
	next http.Handler
}

func (l *loggerConstructor) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	start := time.Now()
	fakerw := responsewriter.NewBuffer(rw)
	defer fakerw.FlushAll()
	l.next.ServeHTTP(fakerw, req)
	end := time.Now()
	latency := end.Sub(start)
	if fakerw.Code != 200 {
		glog.V(1).Infof("\033[33m\033[01m %d \033[0m%s %13v    %s  %s  \033[32m  %s  \033[0m %s  %s", fakerw.Code, req.RemoteAddr, latency, req.Host, req.RequestURI, req.Method, req.Referer(), req.UserAgent())
	} else {
		glog.V(1).Infof("\033[32m %d \033[0m%s %13v    %s  %s  \033[32m %s  \033[0m  %s  %s", fakerw.Code, req.RemoteAddr, latency, req.Host, req.RequestURI, req.Method, req.Referer(), req.UserAgent())
	}
}

// New 日志中间件
func New(ctx context.Context) (alice.Constructor, error) {
	return func(next http.Handler) http.Handler {
		return &loggerConstructor{next: next}
	}, nil
}
