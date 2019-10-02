/*
Copyright 2019 The Labchan Authors.
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

package recovery

import (
	"context"
	"net/http"

	"github.com/golang/glog"
	"github.com/justinas/alice"
)

type recovery struct {
	next http.Handler
}

func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorf("Recovered from panic in http handler: %+v", err)
			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}()
	re.next.ServeHTTP(rw, req)
}

// New creates recovery middleware.
func New(ctx context.Context) (alice.Constructor, error) {
	return func(next http.Handler) http.Handler {
		return &recovery{
			next: next,
		}
	}, nil
}
