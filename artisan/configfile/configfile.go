package configfile

import (
	"github.com/unravela/artisan/api"
	"github.com/unravela/artisan/artisan/configfile/hcl"
	"path/filepath"
)

const (
	// ModuleFile represent name of the file in each module
	ModuleFile = "MODULE.hcl"

	// WorkspaceFile represent name of the root file
	WorkspaceFile = "WORKSPACE.hcl"
)

// LoadModule forward loading of Module into concrete
// implementation. Can be HCL or YAML
func LoadModule(moduleDir string, m *api.Module) error {
	path := filepath.Join(moduleDir, ModuleFile)
	return hcl.LoadModule(path, m)
}

// LoadWorkspace forward loading of Workspace into
// concrete implementation. Can be HCL or YAML.
func LoadWorkspace(workspaceDir string, ws *api.Workspace) error {
	path := filepath.Join(workspaceDir, WorkspaceFile)
	return hcl.LoadWorkspace(path, ws)
}
