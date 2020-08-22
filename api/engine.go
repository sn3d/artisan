package api

// Engine ...
type Engine struct {
	Registry ImageRegistry
	Executor TaskExecutor
}

// Run is main function that runs tasks.
/*
func (r *Runner) Run(task Ref, ws *Workspace) error {
	t, err := ws.Task(task)
	if err != nil {
		return err
	}
	allTasks := topoSort(t, ws)

	for _, tsk := range allTasks {
		err := r.runSingle(tsk, ws)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Runner) runSingle(t *Task, ws *Workspace) error {
	class, _ := ws.Class(t.Class)
	image, err := r.Builder.Build(class, ws)
	if err != nil {
		return err
	}

	err = r.Executor.Exec(t, image, ws.rootDir)
	if err != nil {
		return err
	}

	return nil
}*/
