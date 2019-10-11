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

package cors

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/justinas/alice"
	"tracfox.io/tracfox/internal/util"
)

type corsFilter struct {
	allowOrigin      string
	allowMethods     string
	allowHeaders     string
	allowCredentials bool
	controlMaxAge    int
}

type cors struct {
	plug *corsFilter
	next http.Handler
}

func (c *cors) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Access-Control-Allow-Origin", c.plug.allowOrigin)
	rw.Header().Set("Access-Control-Allow-Methods", c.plug.allowMethods)
	rw.Header().Set("Access-Control-Allow-Headers", c.plug.allowHeaders)
	rw.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(c.plug.allowCredentials))
	rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(c.plug.controlMaxAge))
	if req.Method == "OPTIONS" {
		rw.WriteHeader(http.StatusNoContent)
		return
	}
	c.next.ServeHTTP(rw, req)
}

// Constructor .
func Constructor(args map[string]interface{}) (alice.Constructor, error) {
	var errs []error
	_, ok := args["allowOrigin"]
	if !ok {
		errs = append(errs, errors.New("cors filter allowOrigin field is required"))
	}

	allowOrigin, ok := args["allowOrigin"].(string)
	if !ok {
		errs = append(errs, errors.New("cors filter allowOrigin must be string"))
	}

	_, ok = args["allowMethods"]
	if !ok {
		errs = append(errs, errors.New("cors filter allowMethods field is required"))
	}
	allowMethods, ok := args["allowMethods"].(string)
	if !ok {
		errs = append(errs, errors.New("cors filter allowMethods must be string"))
	}

	_, ok = args["allowHeaders"]
	if !ok {
		errs = append(errs, errors.New("cors filter allowHeaders field is required"))
	}
	allowHeaders, ok := args["allowHeaders"].(string)
	if !ok {
		errs = append(errs, errors.New("cors filter allowHeaders must be string"))
	}

	allowCredentials := false
	_, ok = args["allowCredentials"]
	if ok {
		allowCredentials, ok = args["allowCredentials"].(bool)
		if !ok {
			errs = append(errs, errors.New("cors filter allowCredentials must be bool"))
		}
	}

	controlMaxAge := 0
	_, ok = args["controlMaxAge"]
	if ok {
		controlMaxAge, ok = args["controlMaxAge"].(int)
		if !ok {
			errs = append(errs, errors.New("cors filter controlMaxAge must be int"))
		}
	}

	if len(errs) != 0 {
		return nil, util.NewAggregateError(errs)
	}
	return func(next http.Handler) http.Handler {
		return &cors{
			next: next,
			plug: &corsFilter{
				allowOrigin:      allowOrigin,
				allowMethods:     allowMethods,
				allowHeaders:     allowHeaders,
				allowCredentials: allowCredentials,
				controlMaxAge:    controlMaxAge,
			},
		}
	}, nil
}
