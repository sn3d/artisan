package main

import (
	"github.com/unravela/artisan/cmd/run"
	"github.com/urfave/cli/v2"
	"os"
)

// version is set by goreleaser, via -ldflags="-X 'main.version=...'".
var version = "development"

func main() {
	app := &cli.App{
		Name:  "artisan",
		Version: version,
		Usage: "Artisan build orchestrator tool",
		Commands: []*cli.Command{
			run.RunCmd,
		},
	}

	app.Run(os.Args)
}
