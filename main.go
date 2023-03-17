package main

import (
	"github.com/go-zoox/cli"
)

func main() {
	app := cli.NewSingleProgram(&cli.SingleProgramConfig{
		Name:    "gomirror-hack-gitlab",
		Usage:   "Go mirror hack gitlab preflight => https://stackoverflow.com/questions/56938451/using-go-modules-with-private-repositories-via-http-without-ssl",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "source",
				Aliases:  []string{"s"},
				EnvVars:  []string{"SOURCE"},
				Required: true,
			},
			&cli.StringFlag{
				Name:     "target",
				Aliases:  []string{"t"},
				EnvVars:  []string{"TARGET"},
				Required: true,
			},
		},
	})

	app.Command(func(ctx *cli.Context) error {
		return Serve(&Config{
			Source: ctx.String("source"),
			Target: ctx.String("target"),
		})
	})

	app.Run()
}
