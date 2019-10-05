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
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"

	"tracfox.io/tracfox/internal/responsewriter"
	"tracfox.io/tracfox/internal/util"
)

// DefaultMaxRoundRobinRetryMultiple .
const DefaultMaxRoundRobinRetryMultiple = 3

type roundRobinHandler struct {
	cfg          *Config
	lock         sync.RWMutex
	next         http.Handler
	failureTimes int
	lastFailTime int64
	failMode     bool
}

func (r *roundRobinHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fakeResponseWriter := responsewriter.NewBuffer(rw)
	defer fakeResponseWriter.FlushAll()
	r.next.ServeHTTP(fakeResponseWriter, req)
	r.lock.RLock()

	if util.IsSuccessCode(fakeResponseWriter.Code) && !r.failMode { // success pass
		r.lock.RUnlock()
		return
	}

	if util.IsSuccessCode(fakeResponseWriter.Code) {
		r.lock.RUnlock()
		r.lock.Lock()
		defer r.lock.Unlock()
		r.failMode = false
		r.failureTimes = 0
		return
	}

	r.lock.RUnlock()
	r.lock.Lock()
	defer r.lock.Unlock()
	r.failureTimes++
	r.lastFailTime = time.Now().UnixNano()
	if r.failureTimes > r.cfg.Maxfails {
		r.failMode = true
	}
}

// RoundRobinBalancer .
type RoundRobinBalancer struct {
	weightSum int
	handlers  []*roundRobinHandler
}

// Elect .
func (r *RoundRobinBalancer) Elect(req *http.Request) (http.Handler, error) {
	lenHandlers := len(r.handlers)
	if lenHandlers == 0 {
		return nil, ErrEmptyBackendList
	}
	for i := 0; i < lenHandlers*DefaultMaxRoundRobinRetryMultiple; i++ {
		sum := 0
		if r.weightSum == 0 {
			handler := r.handlers[rand.Intn(lenHandlers)]
			if !handler.failMode || time.Now().UnixNano()-handler.lastFailTime > handler.cfg.FailTimeout {
				return handler, nil
			}
		} else {
			randNum := rand.Intn(r.weightSum)
			for i := 0; i < lenHandlers; i++ {
				sum += r.handlers[i].cfg.Weight
				if sum > randNum {
					if !r.handlers[i].failMode || time.Now().UnixNano()-r.handlers[i].lastFailTime > r.handlers[i].cfg.FailTimeout {
						return r.handlers[i], nil
					}
				}
			}
		}
	}
	return r.handlers[rand.Intn(lenHandlers)], nil
}

// Append .
func (r *RoundRobinBalancer) Append(handler http.Handler, cfg *Config) {
	r.handlers = append(r.handlers, &roundRobinHandler{next: handler, cfg: cfg})
	sort.Sort(r)
	r.weightSum += cfg.Weight
}

func (r *RoundRobinBalancer) Len() int {
	return len(r.handlers)
}

func (r *RoundRobinBalancer) Less(i, j int) bool {
	return r.handlers[i].cfg.Weight < r.handlers[j].cfg.Weight
}

func (r *RoundRobinBalancer) Swap(i, j int) {
	r.handlers[i], r.handlers[j] = r.handlers[j], r.handlers[i]
}
