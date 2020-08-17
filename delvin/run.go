package delvin

// RunContext structure is basically context in which are tasks running
type RunContext struct {
	Builder  ClassBuilder
	Executor TaskExecutor
}

// Run is main function that runs tasks.
func (r *RunContext) Run(t *Task, ws *Workspace) error {
	allTasks := topoSort(t, ws)
	for _, tsk := range allTasks {
		err := r.runSingle(tsk, ws)
		if err != nil {
			return err
		}
	}
	return nil
}

// Dependencies are creating graph where each task might have
// edges pointing to another dependency task. Before we run the
// task, we need to sort the dependency tasks and execute them in
// right order.
//
// This implementation is using Kahn Topology Sort alg. and order
// is reverted. Regular sorting returns you 'a,b,c', this
// implementation returns you 'c,b,a'
//
func topoSort(task *Task, ws *Workspace) []*Task {

	// get whole graph of dependencies for 'task'
	topo := map[Ref]*Task{}
	getAllDeps(task, ws, topo)
	topo[task.Ref] = task

	// create indegree map.
	indegree := map[string]int{}
	for _, t := range topo {
		for _, dep := range t.Deps {
			indegree[dep]++
		}
	}

	// create queue where we will place nodes
	// with no edges.
	queue := []*Task{task}

	// into result we will place tasks in right sorted order
	idx := len(topo)
	result := make([]*Task, idx)

	// This is main kahn topology sort loop which will
	// end when queue is empty
	for len(queue) > 0 {
		// pull task from queue
		picked := queue[0]
		queue = queue[1:]

		// remove dependency edges
		for _, dep := range picked.Deps {
			indegree[dep]--
			if indegree[dep] == 0 {
				appt := ws.Task(Ref(dep))
				queue = append(queue, appt)
			}
		}

		// place the picked into result
		idx--
		result[idx] = picked
	}

	return result
}

func getAllDeps(task *Task, ws *Workspace, allDeps map[Ref]*Task) {
	allDeps[task.Ref] = task

	for _, dep := range task.Deps {
		// normalize ref. - enrich path if is't missing
		ref := Ref(dep)
		if ref.GetPath() == "" {
			ref = NewRef(ref.GetWorkspace(), task.Ref.GetPath(), ref.GetTask())
		}

		depTask := ws.Task(ref)
		getAllDeps(depTask, ws, allDeps)
	}

	return
}

func (r *RunContext) runSingle(t *Task, ws *Workspace) error {
	class := ws.Class(t.Class)
	image, err := r.Builder.Build(class, ws)
	if err != nil {
		return err
	}

	buildDir := ws.AbsPath("//")
	err = r.Executor.Exec(t, image, buildDir)
	if err != nil {
		return err
	}

	return nil
}
