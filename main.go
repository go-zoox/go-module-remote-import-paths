package main

import (
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "go-module-remote-import-paths",
		Usage:   "Go module remote import paths => https://stackoverflow.com/questions/56938451/using-go-modules-with-private-repositories-via-http-without-ssl",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "git-server",
				EnvVars:  []string{"GIT_SERVER"},
				Required: true,
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		return Serve(&Config{
			GitServer: ctx.String("git-server"),
		})
	})

	app.Run()
}
