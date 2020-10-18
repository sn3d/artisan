package version

import "github.com/urfave/cli/v2"

var Cmd = &cli.Command{
	Name:      "version",
	Usage:     "print the version",
	Action: func(context *cli.Context) error {
		cli.ShowVersion(context)
		return nil
	},
}

