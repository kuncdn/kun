package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/glog"
	"tracfox.io/tracfox/internal/util"
	"tracfox.io/tracfox/pkg/tracfox/config"
	"tracfox.io/tracfox/pkg/tracfox/middleware"
	"tracfox.io/tracfox/pkg/tracfox/plugin"
	"tracfox.io/tracfox/pkg/tracfox/proxy"
)

// Server .
type Server struct {
	name           string
	cfg            config.TracwayConfiguration
	ssl            bool
	certificate    string
	certificateKey string
	server         *http.Server
}

func describeBackend(c config.TracwayConfiguration, name string) (config.Backend, error) {
	for i := 0; i < len(c.Backends); i++ {
		if c.Backends[i].Name == name {
			return c.Backends[i], nil
		}
	}
	return config.Backend{}, fmt.Errorf("backend %s not found", name)
}

func describeFrontend(c config.TracwayConfiguration, name string) (config.Frontend, error) {
	for i := 0; i < len(c.Frontends); i++ {
		if c.Frontends[i].Name == name {
			return c.Frontends[i], nil
		}
	}
	return config.Frontend{}, fmt.Errorf("frontend %s not found", name)
}

// New .
func New(ctx context.Context, frontName string, cfg config.TracwayConfiguration) (*Server, error) {
	front, err := describeFrontend(cfg, frontName)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}
	rules := make([]*Rule, 0)
	for _, rule := range front.Rules {
		backend, err := describeBackend(cfg, rule.UseBackend)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		balancer, err := proxy.NewBalancer(backend.Balance)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		// 添加后端
		for _, v := range backend.Servers {
			balancer.Append(proxy.NewBalancedReverseProxy(v), &proxy.Config{
				Name:        v.Name,
				Weight:      v.Weight,
				FailTimeout: (time.Duration(v.FailTimeout) * time.Second).Nanoseconds(), // 纳秒
				Maxfails:    v.Maxfails,
			})
		}

		pluginsChain, err := plugin.NewChain(rule.UsePlugins)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}

		rules = append(rules, &Rule{
			name:          rule.Name,
			methods:       rule.MatchMethods,
			rewiteURIPath: rule.RewiteURIPath,
			reg:           regexp.MustCompile(rule.LocationRegexp),
			chain:         pluginsChain,
			balancer:      balancer,
		})
	}

	chain, err := middleware.NewDefaultChain(ctx)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	return &Server{
		name:           frontName,
		cfg:            cfg,
		ssl:            front.SSL,
		certificate:    front.Certificate,
		certificateKey: front.CertificateKey,
		server: &http.Server{
			Handler:           chain.Then(&RuleMgr{rules: rules}),
			Addr:              front.Address,
			ReadTimeout:       time.Duration(cfg.Default.ReadTimeout) * time.Second,
			WriteTimeout:      time.Duration(cfg.Default.WriteTimeout) * time.Second,
			IdleTimeout:       time.Duration(cfg.Default.IdleTimeout) * time.Second,
			MaxHeaderBytes:    cfg.Default.MaxHeaderBytes,
			ReadHeaderTimeout: time.Duration(cfg.Default.ReadHeaderTimeout) * time.Second,
		},
	}, nil
}

// Run .
func (s *Server) Run() error {
	if !s.ssl {
		glog.Infof("Frontend %s ListenAndServe At %s", s.name, s.server.Addr)
		return s.server.ListenAndServe()
	}
	glog.Infof("Frontend %s ListenAndServeTLS At %s", s.name, s.server.Addr)
	return s.server.ListenAndServeTLS(s.certificate, s.certificateKey)
}

// GracefulStop .
func (s *Server) GracefulStop(ctx context.Context) error {
	glog.Infof("Stoping Frontend %s...", s.name)
	return s.server.Shutdown(ctx)
}

// Manager .
type Manager struct {
	servers []*Server
}

// NewManager .
func NewManager(ctx context.Context, cfg config.TracwayConfiguration) (*Manager, error) {
	servers := make([]*Server, 0)
	for _, frontend := range cfg.Frontends {
		server, err := New(ctx, frontend.Name, cfg)
		if err != nil {
			glog.Exitln(err)
			continue
		}
		servers = append(servers, server)
	}
	return &Manager{servers: servers}, nil
}

// Run .
func (m *Manager) Run() error {
	for _, serv := range m.servers {
		go func(serv *Server) {
			if err := serv.Run(); err != nil {
				glog.Exitln(err)
			}
		}(serv)
	}
	return nil
}

// GracefulStop .
func (m *Manager) GracefulStop(ctx context.Context) error {
	errs := make([]error, 0)
	for _, serv := range m.servers {
		if err := serv.GracefulStop(ctx); err != nil {
			glog.Errorln(err)
			errs = append(errs, err)
		}
	}
	return util.NewAggregateError(errs)
}
