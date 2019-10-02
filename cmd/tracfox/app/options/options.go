package options

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

import (
	"errors"

	"github.com/spf13/pflag"
)

const defaultConfigFile = "/etc/tracfox/config.yaml"

// TracwayFlags 包含帐户的命令行参数。
//如果满足以下任何条件，配置字段应该在TracwayFlags而不是TracwayConfiguration中：
//  - 在节点的生命周期内，它的值永远不会或不能安全地更改，或者
//  - 它的值不能同时在节点之间安全地共享（例如主机名）;
//AccountConfiguration旨在在节点之间共享。
//一般情况下，请尽量避免添加标记或配置字段，
//因为我们已经有了大量令人困惑的东西。
type TracwayFlags struct {
	TracwayConfig string
	DryRun        bool
}

// NewTracwayFlags 将会创建一个新的 TracwayFlags结构，并且填充默认值
func NewTracwayFlags() *TracwayFlags {
	return &TracwayFlags{
		TracwayConfig: defaultConfigFile,
	}
}

// ValidateTracwayFlags 验证TracwayFlags 中填充的数值是否满足要求
func ValidateTracwayFlags(f *TracwayFlags) (errs []error) {
	if len(f.TracwayConfig) == 0 {
		errs = append(errs, errors.New("configuration path is required"))
	}
	return errs
}

// AddFlags adds flags for a specific AccountFlags to the specified FlagSet
func (f *TracwayFlags) AddFlags(mainfs *pflag.FlagSet) {
	fs := pflag.NewFlagSet("", pflag.ExitOnError)
	defer func() {
		mainfs.AddFlagSet(fs)
	}()
	fs.StringVar(&f.TracwayConfig, "config", f.TracwayConfig, "The Tracfox Server will load its initial configuration from this file. The path may be absolute or relative; relative paths start at the Tracfox's current working directory. Omit this flag to use the built-in default configuration values. Command-line flags override configuration from this file.")
	fs.BoolVar(&f.DryRun, "dry-run", f.DryRun, "If true, only check the configuration file and exit.")
}
