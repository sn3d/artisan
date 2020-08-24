package main

import (
	"github.com/unravela/artisan/cmd/run"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "artisan",
		Usage: "Artisan build orchestrator tool",
		Commands: []*cli.Command{
			run.RunCmd,
		},
	}

	app.Run(os.Args)
}
