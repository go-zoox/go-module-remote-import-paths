package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-zoox/core-utils/cast"
	"github.com/go-zoox/proxy"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Config struct {
	// http://10.0.0.1:8888
	GitServer string `json:"git_server"`
	//
	GitServerMap map[string]string `json:"git_server_map"`
	//
	EnableProxy bool `json:"enable_proxy"`
	// custom domain
	RootURL string `json:"root_url"`
}

func Serve(cfg *Config) error {
	if cfg.EnableProxy && cfg.RootURL == "" {
		return fmt.Errorf("RootURL is required when enable proxy")
	}

	app := defaults.Application()

	routes := []proxy.MultiHostsRoute{}
	for host, git := range cfg.GitServerMap {
		u, _ := url.Parse(git)
		ps := u.Port()
		port := int64(80)
		if ps != "" {
			port = cast.ToInt64(ps)
		} else {
			if u.Scheme == "https" {
				port = 443
			}
		}

		routes = append(routes, proxy.MultiHostsRoute{
			Host: host,
			Backend: proxy.MultiHostsRouteBackend{
				ServiceProtocol: u.Scheme,
				ServiceName:     u.Hostname(),
				ServicePort:     port,
			},
		})
	}
	py := proxy.NewMultiHosts(&proxy.MultiHostsConfig{
		Routes: routes,
	})

	app.Get("/*", func(ctx *zoox.Context) {
		if ctx.Query().Get("go-get").String() != "1" {
			if cfg.EnableProxy {
				py.ServeHTTP(ctx.Writer, ctx.Request)
			}
			return
		}

		host := ctx.Host()
		path := strings.TrimSuffix(ctx.Path, "/")
		gitServer := cfg.GitServer
		if cfg.EnableProxy {
			gitServer = cfg.RootURL
		}

		if cfg.GitServerMap != nil {
			if v, ok := cfg.GitServerMap[host]; ok {
				gitServer = v
			}
		}

		importPrefix := host + path
		repoRoot := fmt.Sprintf("%s%s", gitServer, path)

		ctx.HTML(200, BuildGoImport(importPrefix, repoRoot))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return app.Run(fmt.Sprintf(":%s", port))
}
