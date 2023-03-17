package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-zoox/proxy"
	"github.com/go-zoox/proxy/utils/rewriter"
	"github.com/go-zoox/zoox"
	"github.com/go-zoox/zoox/defaults"
)

type Config struct {
	// http://10.0.0.1:8888
	Source string `json:"source"`
	// http://xxx.idp.example.com
	Target string `json:"source"`
}

func (c *Config) SourceHostPort() string {
	u, err := url.Parse(c.Source)
	if err != nil {
		panic("invalid source")
	}

	return u.Host
}

func (c *Config) TargetHostPort() string {
	u, err := url.Parse(c.Target)
	if err != nil {
		panic("invalid target")
	}

	return u.Host
}

func Serve(cfg *Config) error {
	app := defaults.Application()

	// app.Use(func(ctx *zoox.Context) {
	// 	fmt.Println("ctx.Origin():", ctx.Origin())
	// 	fmt.Println("ctx.Host():", ctx.Host())

	// 	ctx.Next()
	// })

	// app.Use(middleware.Proxy(&middleware.ProxyConfig{
	// 	Rewrites: middleware.ProxyGroupRewrites{
	// 		{
	// 			Name:   "gitlab",
	// 			RegExp: "/(.*)",
	// 			Rewrite: middleware.ProxyRewrite{
	// 				upstream: os.Getenv("upstream"),
	// 				Rewrites: rewriter.Rewriters{
	// 					{
	// 						From: "/(.*)",
	// 						To:   "/$1",
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }))

	app.Use(zoox.WrapH(proxy.NewSingleTarget(cfg.Source, &proxy.SingleTargetConfig{
		ChangeOrigin: true,
		Rewrites: rewriter.Rewriters{
			{
				From: "/(.*)",
				To:   "/$1",
			},
		},
		OnResponse: func(resp *http.Response) error {
			goGet := resp.Request.URL.Query().Get("go-get")
			if goGet != "1" {
				return nil
			}

			b, err := ioutil.ReadAll(resp.Body) //Read html
			if err != nil {
				return err
			}
			err = resp.Body.Close()
			if err != nil {
				return err
			}
			bodyString := string(b)
			bodyStringNew := strings.ReplaceAll(bodyString, cfg.SourceHostPort(), cfg.TargetHostPort())
			body := ioutil.NopCloser(bytes.NewReader([]byte(bodyStringNew)))
			resp.Body = body
			resp.ContentLength = int64(len(bodyStringNew))
			resp.Header.Set("Content-Length", strconv.Itoa(len(bodyStringNew)))

			return nil
		},
	})))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return app.Run(fmt.Sprintf(":%s", port))
}
