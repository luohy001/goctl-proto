package main

import (
	"context"
	"fmt"
	"github.com/urfave/cli/v3"
	"runtime"
)

type app struct{}

func (app *app) Run(ctx context.Context, args []string) error {
	var cliApp = cli.Command{}
	cliApp.Name = "goctl-proto"
	cliApp.Usage = "go-zero api file -> proto file"
	cliApp.Version = fmt.Sprintf("%s %s/%s build on %s", max(buildVersion, "v1.0.3"), runtime.GOOS, runtime.GOARCH, max(buildTime, "2024-01-19T14:38:20"))
	cliApp.Commands = []*cli.Command{
		{
			Name:  "proto",
			Usage: "generate proto file from api file",
			UsageText: `goctl-proto proto --input ./example/api/service.api --output ./example
OR with goctl
goctl api plugin -plugin goctl-proto="proto" -api ./example/api/service.api -dir ./example`,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "input",
					Aliases: []string{"i"},
					Usage:   "input api file path",
				},
				&cli.StringFlag{
					Name:    "output",
					Aliases: []string{"o"},
					Usage:   "output proto file directory",
				},
				&cli.StringFlag{
					Name:    "base",
					Aliases: []string{"b"},
					Usage:   "base.proto",
				},
				&cli.StringSliceFlag{
					Name:    "include-handler",
					Aliases: []string{"inc"},
					Usage:   "include handler in api file, prior to exclude-handler",
				},
				&cli.StringSliceFlag{
					Name:    "exclude-handler",
					Aliases: []string{"exc"},
					Usage:   "exclude handler in api file",
				},
				&cli.BoolFlag{
					Name:    "multiple",
					Aliases: []string{"m"},
					Usage:   "output multiple service by api server group",
				},
			},
			Action: protoGen,
		},
	}

	return cliApp.Run(ctx, args)
}
