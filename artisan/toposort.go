package artisan

import (
	"fmt"
	"github.com/unravela/artisan/api"
)

// Dependencies are creating graph where each task might have
// edges pointing to another dependency task. Before we run the
// task, we need to sort the dependency tasks and execute them in
// right order.
//
// This implementation is using Kahn Topology Sort alg. and order
// is reverted. Regular sorting returns you 'a,b,c', this
// implementation returns you 'c,b,a'
//
func topoSort(task *api.Task, inst *Artisan) (api.Tasks, error) {

	// get whole graph of dependencies for 'task'
	topo := map[api.Ref]*api.Task{}
	err := getAllDeps(task, inst, topo)
	if err != nil {
		return nil, err
	}
	topo[task.Ref] = task

	indegree := map[api.Ref]int{}
	for _, t := range topo {
		deps := t.GetDeps()
		for _, dep := range deps {
			indegree[dep]++
		}
	}

	// queue, where we will place nodes
	// with no edges.
	queue := []*api.Task{task}

	// result, where we will place tasks in right sorted order
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
				depRef := dep
				if depRef.IsOnlyTask() {
					depRef = api.NewRef(task.Ref.GetWorkspace(), task.Ref.GetPath(), depRef.GetTask())
				}

				appt, _ := inst.Task(depRef)
				queue = append(queue, appt)
			}
		}

		// place the picked into result
		idx--
		result[idx] = picked
	}

	return result, nil
}

func getAllDeps(task *api.Task, inst *Artisan, allDeps map[api.Ref]*api.Task) error {
	allDeps[task.Ref] = task

	for _, dep := range task.Deps {
		// normalize ref. - enrich path if it's not missing
		ref := api.StringToRef(dep)
		if ref.GetPath() == "" {
			ref = api.NewRef(ref.GetWorkspace(), "//"+task.Ref.GetPath(), ref.GetTask())
		}

		depTask, err := inst.Task(ref)
		if err != nil {
			return fmt.Errorf("Invalid dependency in %s: %w", task.Ref, err)
		}

		err = getAllDeps(depTask, inst, allDeps)
		if err != nil {
			return err
		}
	}

	return nil
}
