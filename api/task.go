package api

// TaskHash is unique value that determine if task need to be re-run
// or everything is up-to-date. It's unique value of state. e.g. If we
// had task 'build' with hash 123, and then hash changed to 321, that
// means some files were changed and probably we should re-run the task.
type TaskHash uint64

// Task determine what to build, what faction to use and specify the dependencies.
type Task struct {
	// FactionName point to faction the task is associated with. Every task must be associated with some faction
	FactionName string `hcl:"faction,label"`
	// Name of the task
	Name string `hcl:"name,label"`
	// Cmd is shell command that is invoked as task
	Cmd []string `hcl:"cmd"`
	// Deps is list of other tasks this task depends on
	Deps []string `hcl:"deps,optional"`

	Include []string `hcl:"include,optional"`
	Exclude []string `hcl:"exclude,optional"`

	// reference to this task.
	Ref Ref
}

// GetDeps returns you valid refs. as dependencies
func (t *Task) GetDeps() []Ref {
	refs := make([]Ref, len(t.Deps))
	for i, dep := range t.Deps {
		// this fragment ensure the ':install' is transformed int '//path/to:install'
		ref := Ref(dep)
		if ref.IsOnlyTask() {
			ref = NewRef(t.Ref.GetWorkspace(), "//" + t.Ref.GetPath(), ref.GetTask())
		}
		refs[i] = ref
	}
	return refs
}