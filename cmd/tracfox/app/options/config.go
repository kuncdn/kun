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

package options

import (
	"crypto/tls"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/pflag"
	"gopkg.in/go-playground/validator.v9"
	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/config"
	"tracfox.io/tracfox/pkg/tracfox/plugin"
	"tracfox.io/tracfox/pkg/tracfox/proxy"
)

// AddTracfoxConfigurationFlags 将config.TracfoxConfiguration对应的所有flag添加到指定的  pflag.FlagSet中
func AddTracfoxConfigurationFlags(mainfs *pflag.FlagSet, f *config.TracfoxConfiguration) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		mainfs.AddFlagSet(fs)
	}()
	fs.StringVar(&f.Default.MetricAddr, "metrics", f.Default.MetricAddr, "Metric address for tracfox server")
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

//ValidateTracfoxConfiguration 验证 TracfoxConfiguration中填充的数据是否满足要求
func ValidateTracfoxConfiguration(f *config.TracfoxConfiguration) (errs []error) {
	if err := validateDefault(f.Default); err != nil {
		errs = append(errs, err)
	}

	if err := validateBackends(f.Backends); err != nil {
		errs = append(errs, err)
	}

	if err := validateFrontends(f.Frontends, f.Backends); err != nil {
		errs = append(errs, err)
	}
	return
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func validateFrontends(c []config.Frontend, b []config.Backend) error {
	errs := make([]error, 0)
	frontNameMap := make(map[string]struct{})
	for _, front := range c {
		if err := validate.Struct(front); err != nil {
			errs = append(errs, err)
		}

		if len(front.Certificate) != 0 || len(front.CertificateKey) != 0 { // check ssl pem
			_, err := tls.LoadX509KeyPair(front.Certificate, front.CertificateKey)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if _, ok := frontNameMap[front.Name]; ok { // check front name
			errs = append(errs, fmt.Errorf("duplicate frontend name %s", front.Name))
		}

		frontNameMap[front.Name] = struct{}{}

		if err := validateVirtualHost(front.VirtualHosts, b); err != nil {
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}

func validateRules(r []config.Rule, b []config.Backend) error {
	errs := make([]error, 0)
	ruleNameMap := make(map[string]struct{})
	for _, rule := range r {
		if err := validate.Struct(rule); err != nil {
			errs = append(errs, err)
		}

		if _, err := regexp.Compile(rule.LocationRegexp); err != nil {
			errs = append(errs, err)
		}

		if !backendExists(b, rule.Backend) {
			errs = append(errs, fmt.Errorf("backend %s not found", rule.Backend))
		}

		if _, ok := ruleNameMap[rule.Name]; ok {
			errs = append(errs, fmt.Errorf("duplicate rule name  %s", rule.Name))
		}

		ruleNameMap[rule.Name] = struct{}{}

		if err := validatePlugins(rule.Plugins); err != nil {
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}

func validateDefault(c config.Default) error {
	return validate.Struct(c)
}

func backendExists(c []config.Backend, name string) bool {
	for i := 0; i < len(c); i++ {
		if c[i].Name == name {
			return true
		}
	}
	return false
}

func validateBackends(b []config.Backend) error {
	errs := make([]error, 0)
	for _, backend := range b {
		if err := validate.Struct(backend); err != nil {
			errs = append(errs, err)
		}

		if _, err := proxy.NewBalancer(backend.Balance); err != nil {
			errs = append(errs, err)
		}

		if err := validateServers(backend.Servers); err != nil {
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}

func validateServers(s []config.Server) error {
	errs := make([]error, 0)
	serverNameMap := make(map[string]struct{})
	for _, server := range s {
		if err := validate.Struct(server); err != nil {
			errs = append(errs, err)
		}
		if _, ok := serverNameMap[server.Name]; ok {
			errs = append(errs, fmt.Errorf("duplicate server name %s", server.Name))
		}
		serverNameMap[server.Name] = struct{}{}
	}
	return util.NewAggregateError(errs)
}

func validatePlugins(ps []config.Plugin) error {
	errs := make([]error, 0)
	for _, plug := range ps { //Plugin error checking
		if err := validate.Struct(plug); err != nil {
			errs = append(errs, err)
		}
		plugConstructor, err := plugin.DescribePlugin(plug.Name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if _, err := plugConstructor(plug.Config); err != nil {
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}

func compileAllRegexp(raws []string) (regs []*regexp.Regexp, err error) {
	for _, raw := range raws {
		reg, err := regexp.Compile(raw)
		if err != nil {
			return nil, err
		}
		regs = append(regs, reg)
	}
	return
}

func validateVirtualHost(vs []config.VirtualHost, b []config.Backend) error {
	errs := make([]error, 0)
	for _, v := range vs {
		if _, err := compileAllRegexp(v.Domains); err != nil {
			errs = append(errs, err)
		}
		if err := validateRules(v.Rules, b); err != nil {
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}
