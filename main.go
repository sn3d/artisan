package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/unravela/artisan/cmd/run"
	ver "github.com/unravela/artisan/cmd/version"

)

// version is set by goreleaser, via -ldflags="-X 'main.version=...'".
var version = "development"

var subcommands = []*cli.Command{
	run.Cmd,
	ver.Cmd,
}

var helpTemplate = `
Artisan is more build orchestrator than a regular build tool

Usage:
	artisan <command> [flags] [arguments...]

Basic Commands:
	run      Run the task

Other Commands:
	version  Print the version information 

Use "artisan help <command>" for more information about a given command.
`

func main() {
	cli.AppHelpTemplate = helpTemplate
	app := &cli.App{
		Name:    "artisan",
		Version: version,
		Commands: subcommands,
		CustomAppHelpTemplate: helpTemplate,
		Action: func(ctx *cli.Context) error {
			return cli.ShowAppHelp(ctx)
		},
	}

	app.Run(os.Args)
}