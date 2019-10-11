/*
Copyright 2019 The Koala Authors.
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
	"errors"
	"math/rand"
	"net/http"
	"reflect"
	"time"
)

var (
	// ErrEmptyBackendList is used when the list of beckends is empty
	ErrEmptyBackendList = errors.New("can not elect backend, Backends empty")
	// ErrCannotElectBackend is used a backend cannot be elected
	ErrCannotElectBackend = errors.New("cant elect backend")
	// ErrUnsupportedAlgorithm is used when an unsupported algorithm is given
	ErrUnsupportedAlgorithm = errors.New("unsupported balancing algorithm")
	typeRegistry            = make(map[string]reflect.Type)
)

// Balancer A Pool is a set of temporary objects that may be individually saved and retrieved.
type Balancer interface {
	Elect(req *http.Request) (http.Handler, error)
	Append(handler http.Handler, cfg *Config)
}

// Config .
type Config struct {
	Name        string
	Weight      int
	FailTimeout int64 // nanoseconds
	Maxfails    int
}

func init() {
	rand.Seed(time.Now().UnixNano())
	typeRegistry["roundrobin"] = reflect.TypeOf(RoundRobinBalancer{})
}

// NewBalancer .
func NewBalancer(balance string) (Balancer, error) {
	alg, ok := typeRegistry[balance]
	if !ok {
		return nil, ErrUnsupportedAlgorithm
	}
	return reflect.New(alg).Elem().Addr().Interface().(Balancer), nil
}
