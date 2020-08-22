package main

import (
	"github.com/unravela/delvin/cmd/run"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "delvin",
		Usage: "Delvin build orchestrator tool",
		Commands: []*cli.Command{
			run.RunCmd,
		},
	}

	app.Run(os.Args)
}
