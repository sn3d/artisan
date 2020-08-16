package delvin

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"io/ioutil"
	"path/filepath"
)

// WorkspaceFile is concrete name of the module HCL file
const ModuleFile = "module.delvin"

// Module basic entity in workspace. It's analogy to module in Maven,
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

// read the HCL file of module.
// Function get the module's absolute path from rootDir and given module's ref.
//
// Example how to load HCL data for '//apps/webapp' module:
//
//    m := &Module{Ref: "//apps/webapp"}
//    loadModuleFile("./", m)
//
func loadModuleHcl(rootDir string, m *Module) error {
	path := m.Ref.AbsPath(rootDir)
	file := filepath.Join(path, ModuleFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	err = hclsimple.Decode(file+".hcl", data, nil, m)
	if err != nil {
		return fmt.Errorf("cannot read data from module file reason: %w", err)
	}

	// set correct ref. to each task
	for _, t := range m.Tasks {
		t.Ref = m.Ref.SetTask(t.Name)
	}

	return nil
}
