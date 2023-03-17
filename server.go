package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Config struct {
	// http://10.0.0.1:8888
	GitServer string `json:"git_server"`
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	app.Get("/*", func(ctx *zoox.Context) {
		if ctx.Query().Get("go-get").String() != "1" {
			return
		}

		host := ctx.Host()
		path := strings.TrimSuffix(ctx.Path, "/")

		importPrefix := host + path
		repoRoot := fmt.Sprintf("%s%s", cfg.GitServer, path)

		ctx.HTML(200, BuildGoImport(importPrefix, repoRoot))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return app.Run(fmt.Sprintf(":%s", port))
}
