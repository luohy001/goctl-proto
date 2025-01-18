package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/luohy001/goctl-proto/internal/proto"
	"github.com/urfave/cli/v3"
	"github.com/zeromicro/go-zero/tools/goctl/pkg/parser/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
	"os"
	"path/filepath"
	"strings"
)

func protoGen(ctx context.Context, command *cli.Command) (err error) {
	output := command.String("output")
	defer func() {
		fmt.Print("Generate proto file")
		if output != "" {
			fmt.Printf(" %s", output)
		}
		if err != nil {
			fmt.Printf(" [FAILED]\n")
		} else {
			fmt.Printf(" [OK]\n")
		}
	}()
	var goctlPlugin plugin.Plugin
	if goctlPlugin.ApiFilePath = command.String("input"); goctlPlugin.ApiFilePath != "" {
		if goctlPlugin.Api, err = parser.Parse(goctlPlugin.ApiFilePath, ""); err != nil {
			return err
		}
	} else if plug, err := newPlugin(); err == nil {
		goctlPlugin = *plug
	} else {
		return errors.New("api file not found, must set one of goctl -api or --input")
	}
	apiFileName := ""
	if apiFile := filepath.Base(goctlPlugin.ApiFilePath); goctlPlugin.Dir != "" {
		output = filepath.Join(goctlPlugin.Dir, strings.TrimSuffix(apiFile, filepath.Ext(apiFile))+".proto")
		apiFileName = strings.TrimSuffix(apiFile, filepath.Ext(apiFile))
	} else {
		fi, err := os.Stat(output)
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return errors.New("output is not a directory")
		}
		output = filepath.Join(output, strings.TrimSuffix(apiFile, filepath.Ext(apiFile))+".proto")
		apiFileName = strings.TrimSuffix(apiFile, filepath.Ext(apiFile))
	}
	pf, err := proto.Unmarshal(goctlPlugin.Api, command.Bool("multiple"), apiFileName)
	if err != nil {
		return err
	}
	pd, err := pf.Refine(command.StringSlice("include-handler"), command.StringSlice("exclude-handler")).Marshal()
	if err != nil {
		return err
	}
	if err = os.WriteFile(output, pd, 0666); err != nil {
		return err
	}
	return
}

func newPlugin() (*plugin.Plugin, error) {
	if stat, err := os.Stdin.Stat(); err != nil {
		return nil, err
	} else if stat.Size() <= 0 {
		return nil, errors.New("empty stdin")
	}
	return plugin.NewPlugin()
}
