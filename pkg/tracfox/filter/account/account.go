package account

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

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/api"
)

type cors struct {
	next       http.Handler
	clientPool *sync.Pool
}

func (c *cors) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	client := c.clientPool.Get().(api.AccountClient)
	defer c.clientPool.Put(client)
	token := ""
	rawTokens := strings.Split(req.Header.Get("Authorization"), " ")
	if len(rawTokens) == 2 {
		token = rawTokens[1]
	}
	resp, err := client.ValidateUserPermission(req.Context(), &api.PermissionReq{
		Route:  req.URL.Path,
		Token:  token,
		Method: req.Method,
	})
	if err != nil {
		glog.Errorln(err)
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Println(resp.GetRetCode())
	switch resp.GetRetCode() {
	case api.RetCodeType_NO_MATCHED_API:
		http.Error(rw, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	case api.RetCodeType_LOGIN_REQUIRED:
		http.Error(rw, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	case api.RetCodeType_PERMISSION_DENIED:
		http.Error(rw, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	case api.RetCodeType_ALLOW_ACCESS:
		c.next.ServeHTTP(rw, req)
		return
	}
}

// Constructor .
func Constructor(args map[string]interface{}) (alice.Constructor, error) {
	var errs []error
	_, ok := args["serverName"]
	if !ok {
		errs = append(errs, errors.New("labchan filter serverName field is required"))
	}

	serverName, ok := args["serverName"].(string)
	if !ok {
		errs = append(errs, errors.New("labchan filter serverName must be string"))
	}

	_, ok = args["address"]
	if !ok {
		errs = append(errs, errors.New("labchan filter address field is required"))
	}
	address, ok := args["address"].(string)
	if !ok {
		errs = append(errs, errors.New("labchan filter address must be string"))
	}

	_, ok = args["certificate"]
	if !ok {
		errs = append(errs, errors.New("labchan filter certificate field is required"))
	}
	certificate, ok := args["certificate"].(string)
	if !ok {
		errs = append(errs, errors.New("labchan filter certificate must be string"))
	}

	creds, err := credentials.NewClientTLSFromFile(certificate, serverName)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return nil, util.NewAggregateError(errs)
	}

	return func(next http.Handler) http.Handler {
		return &cors{
			next: next,
			clientPool: &sync.Pool{
				New: func() interface{} {
					conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))
					if err != nil {
						glog.Errorln(err.Error())
						return nil
					}
					return api.NewAccountClient(conn)
				},
			},
		}
	}, nil
}
