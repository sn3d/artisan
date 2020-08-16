package delvin

// Task determine what to build, what forge to use and specify the dependencies. Tasks are organized in
// graph.
type Task struct {
	// Class the task belong to. Every task must be associated with some class
	Class string `hcl:"class,label"`
	// Name of the task
	Name string `hcl:"name,label"`
	// Cmd is shell command that is invoked as task
	Cmd []string `hcl:"cmd"`
	// Deps is list of other tasks this task depends on
	Deps []string `hcl:"deps,optional"`

	Include []string `hcl:"inclube,optional"`
	Exclude []string `hcl:"exclude,optional"`

	// reference to this task.
	Ref Ref
}
