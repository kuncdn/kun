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

package config

// KoalaConfiguration .
type KoalaConfiguration struct {
	Default   Default    `yaml:"default" validate:"required"`
	Frontends []Frontend `yaml:"frontends" validate:"required"`
	Backends  []Backend  `yaml:"backends" validate:"required"`
}

// Backend .
type Backend struct {
	Name    string   `yaml:"name" validate:"required"`
	Balance string   `yaml:"balance" validate:"required"`
	Servers []Server `yaml:"servers" validate:"required"`
}

// Server .
type Server struct {
	Name                string `yaml:"name" validate:"required"`
	Weight              int    `yaml:"weight" validate:"min=0"`
	FailTimeout         int    `yaml:"failTimeout" validate:"min=0"`
	MaxIdleConnsPerHost int    `yaml:"maxIdleConnsPerHost" validate:"min=0"`
	Maxfails            int    `yaml:"maxFails" validate:"min=0"`
	TCPTimeout          int    `yaml:"tcpTimeout" validate:"min=0"`
	TCPKeepAlive        int    `yaml:"tcpKeepAlive" validate:"min=0"`
	IdleConnTimeout     int    `yaml:"idleConnTimeout" validate:"min=0"`
	Target              string `yaml:"target" validate:"required"`
}

// Default .
type Default struct {
	ReadTimeout       int `yaml:"readTimeout" validate:"required,min=0"`
	IdleTimeout       int `yaml:"idleTimeout" validate:"required,min=0"`
	WriteTimeout      int `yaml:"writeTimeout" validate:"required,min=0"`
	MaxHeaderBytes    int `yaml:"maxHeaderBytes" validate:"required,min=0"`
	ReadHeaderTimeout int `yaml:"readHeaderTimeout" validate:"required,min=0"`
	GraceTimeOut      int `yaml:"graceTimeOut" validate:"required,min=0"`
}

// Frontend .
type Frontend struct {
	Name           string        `yaml:"name" validate:"required"`
	Address        string        `yaml:"address" validate:"required"`
	Certificate    string        `yaml:"certificate"`
	CertificateKey string        `yaml:"certificateKey"`
	VirtualHosts   []VirtualHost `yaml:"virtualHosts" validate:"required"`
}

// Rule .
type Rule struct {
	Name           string   `yaml:"name" validate:"required"`
	LocationRegexp string   `yaml:"locationRegexp" validate:"required"`
	MatchMethods   []string `yaml:"matchMethods" validate:"required"`
	RewitePath     string   `yaml:"rewitePath"`
	Backend        string   `yaml:"backend" validate:"required"`
	Filters        []Filter `yaml:"filters"`
}

// Filter .
type Filter struct {
	Name   string                 `yaml:"name" validate:"required"`
	Config map[string]interface{} `yaml:"config"`
}

// VirtualHost .
type VirtualHost struct {
	Filters []Filter `yaml:"filters"`
	Domains []string `yaml:"domains" validate:"required"`
	Rules   []Rule   `yaml:"rules" validate:"required"`
}
