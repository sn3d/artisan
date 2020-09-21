package api


// Module basic entity in artisan. It's analogy to module in Maven,
// package in Bazel etc.
type Module struct {
	// Tasks is collection of all tasks the module contain
	Tasks []*Task `hcl:"task,block"`

	// Ref is reference to the module
	Ref Ref
}

// Task gives you task for given name or nil if there
// is no task with this name available
func (m *Module) Task(name string) *Task {
	for _, t := range m.Tasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}
