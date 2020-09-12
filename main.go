package main

import (
	"os"

	"github.com/unravela/artisan/cmd/run"
	"github.com/urfave/cli/v2"
)

// version is set by goreleaser, via -ldflags="-X 'main.version=...'".
var version = "development"

func main() {
	app := &cli.App{
		Name:    "artisan",
		Version: version,
		Usage:   "Artisan build orchestrator",
		Commands: []*cli.Command{
			run.RunCmd,
		},
	}

	app.Run(os.Args)
}
