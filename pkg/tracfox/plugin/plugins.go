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

package plugin

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/justinas/alice"
	"tracfox.io/tracfox/pkg/tracfox/config"
	"tracfox.io/tracfox/pkg/tracfox/plugin/cors"
	"tracfox.io/tracfox/pkg/tracfox/plugin/labchan"
)

func init() {
	plugins["cors"] = cors.Constructor
	plugins["accessByAccount"] = labchan.Constructor
}

var (
	// plugins is a map of plugin name to Plugin.
	plugins = make(map[string]Constructor)
)

// Constructor defines basic methods for plugins
type Constructor func(cfg map[string]interface{}) (alice.Constructor, error)

// RegisterPlugin plugs in plugin. All plugins should register
// themselves, even if they do not perform an action associated
// with a directive. It is important for the process to know
// which plugins are available.
//
// The plugin MUST have a name: lower case and one word.
// If this plugin has an action, it must be the name of
// the directive that invokes it. A name is always required
// and must be unique for the server type.
func RegisterPlugin(name string, plugin Constructor) error {
	if name == "" {
		return errors.New("plugin must have a name")
	}
	if _, dup := plugins[name]; dup {
		return fmt.Errorf("plugin named %s  already registered", name)
	}
	plugins[name] = plugin
	return nil
}

// DescribePlugin gets the action for a plugin
func DescribePlugin(name string) (Constructor, error) {
	if plugin, ok := plugins[name]; ok {
		return plugin, nil
	}
	return nil, fmt.Errorf("plugin %q not found", name)
}

// NewChain .
func NewChain(ps []config.Plugin) (alice.Chain, error) {
	chain := alice.New()
	for _, v := range ps {
		constructor, err := DescribePlugin(v.Name)
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
