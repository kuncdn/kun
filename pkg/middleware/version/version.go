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

package version

import (
	"context"
	"net/http"

	"github.com/justinas/alice"
	uuid "github.com/satori/go.uuid"
)

const serverName = "Koala"

type headers struct {
	next http.Handler
}

func (h *headers) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Server", serverName)         // set server name for response
	req.Header.Set("uuid", uuid.NewV4().String()) // generate request uuid
	h.next.ServeHTTP(rw, req)
}

// New create heade middleware
func New(ctx context.Context) (alice.Constructor, error) {
	return func(next http.Handler) http.Handler {
		return &headers{
			next: next,
		}
	}, nil
}
