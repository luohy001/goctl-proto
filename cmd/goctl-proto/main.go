package main

import (
	"context"
	"fmt"
	"os"
)

var (
	buildTime    string
	buildVersion string
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := new(app).Run(ctx, os.Args); err != nil {
		fmt.Println(err)
	}
}
