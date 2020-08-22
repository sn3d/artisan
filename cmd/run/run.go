package run

import (
	"github.com/unravela/delvin/api"
	"github.com/unravela/delvin/workspace"
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
	Action:    RunAction,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "workspace-dir",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "Path to workspace directory, by default it's current directory.",
			Destination: &opts.CurrentDir,
		},
	},
}

func RunAction(ctx *cli.Context) error {
	return Run(ctx.Args().First(), opts)
}

// Run the given task
func Run(task string, opts RunOpts) error {
	ws, err := workspace.Open(opts.CurrentDir)
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
