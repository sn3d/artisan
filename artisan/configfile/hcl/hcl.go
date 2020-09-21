package hcl

import (
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/unravela/artisan/api"
)

const (
	// ModuleFile represent name of the file in each module
	ModuleFile = "MODULE.hcl"

	// WorkspaceFile represent name of the root file
	WorkspaceFile = "WORKSPACE.hcl"
)

// workspaceData is main HCL structure for WORKSPACE.hcl file
type workspaceData struct {

	// hold the forges defined in this artisan
	Factions []*api.Faction `hcl:"faction,block"`

	// hold the main module in artisan file
	MainModule *api.Module `hcl:"module,block"`
}

// LoadModule consumes absolute path to MODULE.hcl and
// fill thee module structure with data
func LoadModule(path string, m *api.Module) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	if err = hclsimple.Decode(path, data, nil, m); err != nil {
		return fmt.Errorf("cannot read data from module file reason: %w", err)
	}

	// set correct ref. for each task
	for _, t := range m.Tasks {
		t.Ref = m.Ref.SetTask(t.Name)
	}

	return nil
}

func LoadWorkspace(path string, ws *api.Workspace) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	var hclData workspaceData
	if err = hclsimple.Decode(path, data, nil, &hclData); err != nil {
		return fmt.Errorf("cannot read data from module file reason: %w", err)
	}

	ws.Factions = api.NewFactions(hclData.Factions)
	ws.MainModule = hclData.MainModule

	return nil
}
