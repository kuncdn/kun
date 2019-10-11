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

package server

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/golang/glog"
	"github.com/shimcdn/koala/internal/util"
	"github.com/shimcdn/koala/pkg/config"
	"github.com/shimcdn/koala/pkg/filter"
	"github.com/shimcdn/koala/pkg/middleware"
	"github.com/shimcdn/koala/pkg/proxy"
)

// Server .
type Server struct {
	name           string
	cfg            config.KoalaConfiguration
	certificate    string
	certificateKey string
	server         *http.Server
}

func describeBackend(c config.KoalaConfiguration, name string) (config.Backend, error) {
	for i := 0; i < len(c.Backends); i++ {
		if c.Backends[i].Name == name {
			return c.Backends[i], nil
		}
	}
	return config.Backend{}, fmt.Errorf("backend %s not found", name)
}

func describeFrontend(c config.KoalaConfiguration, name string) (config.Frontend, error) {
	for i := 0; i < len(c.Frontends); i++ {
		if c.Frontends[i].Name == name {
			return c.Frontends[i], nil
		}
	}
	return config.Frontend{}, fmt.Errorf("frontend %s not found", name)
}

func mustCompileAllRegexp(raws []string) (regs []*regexp.Regexp) {
	for _, raw := range raws {
		reg := regexp.MustCompile(raw)
		regs = append(regs, reg)
	}
	return regs
}

// New .
func New(ctx context.Context, frontName string, cfg config.KoalaConfiguration) (*Server, error) {
	front, err := describeFrontend(cfg, frontName)
	if err != nil {
		glog.Errorln(err)
		return nil, err
	}

	virtualHosts := make([]*VirtualHost, 0)
	for _, host := range front.VirtualHosts {
		virtualHostRegs := mustCompileAllRegexp(host.Domains)
		virtualHostFilterChain, err := filter.NewChain(host.Filters)
		if err != nil {
			glog.Errorln(err)
			return nil, err
		}
		rules := make([]*Rule, 0)
		for _, rule := range host.Rules {
			backend, err := describeBackend(cfg, rule.Backend)
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
				balancer.Append(proxy.NewBalancedReverseProxy(rule, v), &proxy.Config{
					Name:        v.Name,
					Weight:      v.Weight,
					FailTimeout: (time.Duration(v.FailTimeout) * time.Second).Nanoseconds(), // Nanoseconds
					Maxfails:    v.Maxfails,
				})
			}

			filtersChain, err := filter.NewChain(rule.Filters)
			if err != nil {
				glog.Errorln(err)
				return nil, err
			}

			rules = append(rules, &Rule{
				name:     rule.Name,
				methods:  rule.MatchMethods,
				reg:      regexp.MustCompile(rule.LocationRegexp),
				chain:    filtersChain,
				balancer: balancer,
			})
		}
		virtualHosts = append(virtualHosts, &VirtualHost{
			next:            virtualHostFilterChain.Then(&RuleMgr{rules: rules}),
			virtualHostRegs: virtualHostRegs,
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
		certificate:    front.Certificate,
		certificateKey: front.CertificateKey,
		server: &http.Server{
			Handler:           chain.Then(&VirtualHostMgr{virtualHosts: virtualHosts}),
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
	if len(s.certificate) != 0 || len(s.certificateKey) != 0 {
		glog.Infof("Frontend %s ListenAndServeTLS At %s", s.name, s.server.Addr)
		return s.server.ListenAndServeTLS(s.certificate, s.certificateKey)
	}
	glog.Infof("Frontend %s ListenAndServe At %s", s.name, s.server.Addr)
	return s.server.ListenAndServe()
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
func NewManager(ctx context.Context, cfg config.KoalaConfiguration) (*Manager, error) {
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
