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
	lock         sync.Mutex
	next         http.Handler
	failureTimes int
	lastFailTime int64
	failMode     bool
}

func (r *roundRobinHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fakeResponseWriter := responsewriter.NewBuffer(rw)
	defer fakeResponseWriter.FlushAll()
	r.next.ServeHTTP(fakeResponseWriter, req)
	r.lock.Lock()
	defer r.lock.Unlock()
	if util.IsSuccessCode(fakeResponseWriter.Code) {
		r.failMode = false
		r.failureTimes = 0
	} else {
		r.failureTimes++
		r.lastFailTime = time.Now().UnixNano()
		if r.failureTimes > r.cfg.Maxfails {
			r.failMode = true
		}
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
