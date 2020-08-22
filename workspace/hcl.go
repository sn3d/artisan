package workspace

import (
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/unravela/delvin/api"
	"io/ioutil"
)


const (
	ModuleFile = "module.delvin"
	WorkspaceFile = "workspace.delvin"
)

type WorkspaceHCL struct {

	// hold the forges defined in this workspace
	Classes []*api.Class `hcl:"class,block"`

	// hold the main module in workspace file
	MainModule *api.Module `hcl:"module,block"`
}

func LoadModuleFromHCL(path string, m *api.Module) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	if err = hclsimple.Decode(path+".hcl", data, nil, m); err != nil {
		return fmt.Errorf("cannot read data from module file reason: %w", err)
	}

	// set correct ref. for each task
	for _, t := range m.Tasks {
		t.Ref = m.Ref.SetTask(t.Name)
	}

	return nil
}


func loadWorkspaceFromHCL(path string, ws *Workspace) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read module file in %s", path)
	}

	var hclData WorkspaceHCL
	if err = hclsimple.Decode(path+".hcl", data, nil, &hclData); err != nil {
		return fmt.Errorf("cannot read data from module file reason: %w", err)
	}

	ws.classes = arrayToClasses(hclData.Classes)
	ws.mainModule = hclData.MainModule

	return nil
}
