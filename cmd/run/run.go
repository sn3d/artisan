package run

import (
	"github.com/unravela/artisan/api"
	"github.com/unravela/artisan/artisan"
	"github.com/urfave/cli/v2"
	"os"
)

type RunOpts struct {
	CurrentDir string
}

var opts RunOpts

var RunCmd = &cli.Command{
	Name:      "run",
	Usage:     "run the task",
	ArgsUsage: "[//path/to/module:task or :task]",
	Action:    runAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "artisan-dir",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "Path to artisan directory, by default it's current directory.",
			Destination: &opts.CurrentDir,
		},
	},
}

func runAction(ctx *cli.Context) error {
	return Run(ctx.Args().First(), opts)
}

// Run the given task
func Run(task string, opts RunOpts) error {
	ws, err := artisan.OpenWorkspace(opts.CurrentDir)
	if err != nil {
		return err
	}

	// resolve the complete ref to task if we put only ':task'
	taskRef := api.Ref(task)
	if taskRef.IsOnlyTask() {
		currentDir, _ := os.Getwd()
		moduleRef := ws.FindModule(currentDir)
		taskRef = moduleRef.SetTask(taskRef.GetTask())
	}

	err = ws.Run(taskRef)
	if err != nil {
		return err
	}

	return nil
}
