package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "go-module-remote-import-paths",
		Usage:   "Go module remote import paths => https://stackoverflow.com/questions/56938451/using-go-modules-with-private-repositories-via-http-without-ssl",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "git-server",
				EnvVars: []string{"GIT_SERVER"},
			},
			&cli.StringFlag{
				Name:    "git-server-map",
				EnvVars: []string{"GIT_SERVER_MAP"},
			},
			&cli.StringFlag{
				Name:    "root-url",
				EnvVars: []string{"ROOT_URL"},
			},
			&cli.BoolFlag{
				Name:    "enable-proxy",
				EnvVars: []string{"ENABLE_PROXY"},
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		gitServer := ctx.String("git-server-map")
		gitServerMap := map[string]string{}
		if gitServer != "" {
			err := json.Unmarshal([]byte(gitServer), &gitServerMap)
			if err != nil {
				return fmt.Errorf("invalid git-server-map(%s): %s", gitServer, err)
			}
		}

		if ctx.Bool("enable-proxy") && ctx.String("root-url") == "" {
			return fmt.Errorf("root-url is required when enable proxy")
		}

		return Serve(&Config{
			RootURL:      ctx.String("root-url"),
			GitServer:    ctx.String("git-server"),
			GitServerMap: gitServerMap,
			EnableProxy:  ctx.Bool("enable-proxy"),
		})
	})

	app.Run()
}
