package workspace

import "github.com/unravela/delvin/api"

// Dependencies are creating graph where each task might have
// edges pointing to another dependency task. Before we run the
// task, we need to sort the dependency tasks and execute them in
// right order.
//
// This implementation is using Kahn Topology Sort alg. and order
// is reverted. Regular sorting returns you 'a,b,c', this
// implementation returns you 'c,b,a'
//
func topoSort(task *api.Task, ws *Workspace) Tasks {

	// get whole graph of dependencies for 'task'
	topo := map[api.Ref]*api.Task{}
	getAllDeps(task, ws, topo)
	topo[task.Ref] = task

	// create indegree map.
	indegree := map[api.Ref]int{}
	for _, t := range topo {
		deps := t.GetDeps()
		for _, dep := range deps {
			indegree[dep]++
		}
	}

	// create queue where we will place nodes
	// with no edges.
	queue := []*api.Task{task}

	// into result we will place tasks in right sorted order
	idx := len(topo)
	result := make([]*api.Task, idx)

	// This is main kahn topology sort loop which will
	// end when queue is empty
	for len(queue) > 0 {
		// pull task from queue
		picked := queue[0]
		queue = queue[1:]

		// remove dependency edges
		deps := picked.GetDeps()
		for _, dep := range deps {
			indegree[dep]--
			if indegree[dep] == 0 {
				// this fragment ensure the ':install' is transformed int '//path/to:install'
				depRef := api.Ref(dep)
				if depRef.IsOnlyTask() {
					depRef = api.NewRef(task.Ref.GetWorkspace(), task.Ref.GetPath(), depRef.GetTask())
				}

				appt, _ := ws.Task(depRef)
				queue = append(queue, appt)
			}
		}

		// place the picked into result
		idx--
		result[idx] = picked
	}

	return result
}

func getAllDeps(task *api.Task, ws *Workspace, allDeps map[api.Ref]*api.Task) {
	allDeps[task.Ref] = task

	for _, dep := range task.Deps {
		// normalize ref. - enrich path if is't missing
		ref := api.Ref(dep)
		if ref.GetPath() == "" {
			ref = api.NewRef(ref.GetWorkspace(), "//" + task.Ref.GetPath(), ref.GetTask())
		}

		depTask, _ := ws.Task(ref)
		getAllDeps(depTask, ws, allDeps)
	}

	return
}
