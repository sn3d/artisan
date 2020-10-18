package run

import (
	"fmt"
	"github.com/unravela/artisan/api"
	"github.com/unravela/artisan/artisan"
	"github.com/urfave/cli/v2"
	"os"
)

type Opts struct {
	CurrentDir string
}

var Cmd = &cli.Command{
	Name:      "run",
	Usage:     "run the task",
	ArgsUsage: "[//path/to/module:task or :task]",
	Action:    main,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "artisan-dir",
			Aliases:     []string{"d"},
			Value:       "",
			Usage:       "Path to artisan directory, by default it's current directory.",
		},
	},
}

func main(ctx *cli.Context) error {
	task := ctx.Args().First()
	if task == "" {
		return fmt.Errorf("no task specified")
	}

	return Run(task, Opts{
		CurrentDir: ctx.String("artisan-dir"),
	})
}

// Run the given task
func Run(task string, opts Opts) error {
	ws, err := artisan.OpenWorkspace(opts.CurrentDir)
	if err != nil {
		return err
	}

	// resolve the complete ref to task if we put only ':task'
	taskRef := api.StringToRef(task)
	if taskRef.IsOnlyTask() {
		currentDir, _ := os.Getwd()
		moduleRef := ws.FindModule(currentDir)
		taskRef = moduleRef.SetTask(taskRef.GetTask())
	}

	err = ws.Run(taskRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err);
		return err
	}

	return nil
}
