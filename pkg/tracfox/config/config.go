package config

// TracwayConfiguration .
type TracwayConfiguration struct {
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
	MetricAddr        string `yaml:"metricAddr" validate:"required"`
	ReadTimeout       int    `yaml:"readTimeout" validate:"required,min=0"`
	IdleTimeout       int    `yaml:"idleTimeout" validate:"required,min=0"`
	WriteTimeout      int    `yaml:"writeTimeout" validate:"required,min=0"`
	MaxHeaderBytes    int    `yaml:"maxHeaderBytes" validate:"required,min=0"`
	ReadHeaderTimeout int    `yaml:"readHeaderTimeout" validate:"required,min=0"`
	GraceTimeOut      int    `yaml:"graceTimeOut" validate:"required,min=0"`
}

// Frontend .
type Frontend struct {
	Name           string `yaml:"name" validate:"required"`
	Address        string `yaml:"address" validate:"required"`
	SSL            bool   `yaml:"ssl"`
	Certificate    string `yaml:"certificate"`
	CertificateKey string `yaml:"certificateKey"`
	Rules          []Rule `yaml:"rules" validate:"required"`
}

// Rule .
type Rule struct {
	Name           string   `yaml:"name" validate:"required"`
	LocationRegexp string   `yaml:"locationRegexp" validate:"required"`
	MatchMethods   []string `yaml:"matchMethods" validate:"required"`
	RewiteURIPath  string   `yaml:"rewiteUriPath"`
	UseBackend     string   `yaml:"useBackend" validate:"required"`
	UsePlugins     []Plugin `yaml:"usePlugins"`
}

// Plugin .
type Plugin struct {
	Name   string                 `yaml:"name" validate:"required"`
	Config map[string]interface{} `yaml:"config"`
}
