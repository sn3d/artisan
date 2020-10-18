package api

// TaskHash is unique value that determine if task need to be re-run
// or everything is up-to-date. It's unique value of state. e.g. If we
// had task 'build' with hash 123, and then hash changed to 321, that
// means some files were changed and probably we should re-run the task.
type TaskHash uint64

// Task determine what to build, what Environment to use and specify the dependencies.
type Task struct {
	// EnvName point to env. the task is associated with. Every task must
	// be associated with some environment
	EnvName string `hcl:"environment,label"`
	// Name of the task
	Name string `hcl:"name,label"`
	// Script is command that is invoked as task. It's executed as shell-form
	// of CMD
	Script string `hcl:"script"`
	// Deps is list of version tasks this task depends on
	Deps []string `hcl:"deps,optional"`

	Include []string `hcl:"include,optional"`
	Exclude []string `hcl:"exclude,optional"`

	// reference to this task.
	Ref Ref
}

// Tasks is preffered type for collection of Task
type Tasks []*Task


// GetDeps returns you valid refs. as dependencies
func (t *Task) GetDeps() []Ref {
	refs := make([]Ref, len(t.Deps))
	for i, dep := range t.Deps {
		// this fragment ensure the ':install' is transformed int '//path/to:install'
		ref := StringToRef(dep)
		if ref.IsOnlyTask() {
			ref = NewRef(t.Ref.GetWorkspace(), "//" + t.Ref.GetPath(), ref.GetTask())
		}
		refs[i] = ref
	}
	return refs
}
