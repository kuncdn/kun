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

package filter

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"github.com/shimcdn/koala/pkg/config"
	"github.com/shimcdn/koala/pkg/filter/account"
	"github.com/shimcdn/koala/pkg/filter/cors"
)

func init() {
	filters["cors"] = cors.Constructor
	filters["accessByAccount"] = account.Constructor
}

var (
	// filters is a map of filter name to Filter.
	filters = make(map[string]Constructor)
)

// Constructor defines basic methods for filters
type Constructor func(cfg map[string]interface{}) (alice.Constructor, error)

// RegisterFilter plugs in filter. All filters should register
// themselves, even if they do not perform an action associated
// with a directive. It is important for the process to know
// which filters are available.
//
// The filter MUST have a name: lower case and one word.
// If this filter has an action, it must be the name of
// the directive that invokes it. A name is always required
// and must be unique for the server type.
func RegisterFilter(name string, filter Constructor) error {
	if name == "" {
		return errors.New("filter must have a name")
	}
	if _, dup := filters[name]; dup {
		return fmt.Errorf("filter named %s  already registered", name)
	}
	filters[name] = filter
	return nil
}

// DescribeFilter gets the action for a filter
func DescribeFilter(name string) (Constructor, error) {
	if filter, ok := filters[name]; ok {
		return filter, nil
	}
	return nil, fmt.Errorf("filter %q not found", name)
}

// NewChain .
func NewChain(ps []config.Filter) (alice.Chain, error) {
	chain := alice.New()
	for _, v := range ps {
		constructor, err := DescribeFilter(v.Name)
		if err != nil {
			glog.Errorln(err)
			return chain, err
		}
		plug, err := constructor(v.Config)
		if err != nil {
			glog.Errorln(err)
			return chain, err
		}
		chain = chain.Append(plug)
	}
	return chain, nil
}
