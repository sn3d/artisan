package delvin

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"io/ioutil"
	"os"
	"path/filepath"
)

// WorkspaceFile is concrete name of the file
const WorkspaceFile = "workspace.delvin"

type Workspace struct {
	// absolute path to root of the workspace. This is a dir where
	// workspace file is situated
	rootDir string

	// hold the forges defined in this workspace
	Classes []*Class `hcl:"class,block"`
}

// Open opens the workspace. If you're in subdirectory of workspace, the
// function will traverse upside until workspace root is found.
func Open(currentDir string) (*Workspace, error) {

	rootDir, err := findRoot(currentDir)
	if err != nil {
		return nil, errors.New("Cannot open workspace")
	}

	ws := &Workspace{
		rootDir: rootDir,
	}

	err = loadWorkspaceFile(rootDir, ws)
	if err != nil {
		return nil, err
	}

	return ws, nil
}

// Class returns a class for the given name
func (ws *Workspace) Class(name string) *Class {
	for _, c := range ws.Classes {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// Module returns a module for given ref, no matter if
// ref. contains also task.
func (ws *Workspace) Module(ref Ref) (*Module, error) {
	module := &Module{
		Ref: ref,
	}

	err := loadModuleHcl(ws.rootDir, module)
	if err != nil {
		return nil, err
	}

	return module, nil
}

// Task gives you task by given reference. If there is no
// task, the null is returned.
func (ws *Workspace) Task(ref Ref) *Task {
	tskName := ref.GetTask()
	if tskName == "" {
		return nil
	}

	module, err := ws.Module(ref)
	if err != nil {
		return nil
	}

	task := module.Task(tskName)
	if task == nil {
		return nil
	}

	return task
}

// this function traverse upside over directories until we found directory with 'moncici.workspace'
// file. The reason why we need to traverse-up the directories is because
// user might call 'moncici' in any subdirectory of repository.
//
// Function is returning absolute path of root directory or empty string with
// error if there is no workspace file.
func findRoot(dir string) (string, error) {
	currentDir, _ := filepath.Abs(dir)

	for {
		wsFile := filepath.Join(currentDir, WorkspaceFile)
		info, err := os.Stat(wsFile)

		// if we found workspace root
		if err == nil && info.Mode().IsRegular() {
			return currentDir, nil
		}

		// ... or there is no workspace file
		if os.IsNotExist(err) {
			// go to upper/parent directory
			prevDir := currentDir
			currentDir, err = filepath.Abs(currentDir + "/..")

			// this check is because when we're on top (e.g. '/' or 'c:/')
			// we need to break the loop. WIthout this break, the find will enter
			// into infinity loop where currentDir is always root dir.
			if prevDir == currentDir {
				break
			}
			continue
		}

		// ... or some unknown error is raised
		if err != nil {
			break
		}
	}

	return "", errors.New("Root of woskpace not fount")
}

// function open the workspace file and load data into structure
func loadWorkspaceFile(wsRootDir string, ws *Workspace) error {
	file := filepath.Join(wsRootDir, WorkspaceFile)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read workspace file in %s", wsRootDir)
	}

	err = hclsimple.Decode(file+".hcl", data, nil, ws)
	if err != nil {
		return fmt.Errorf("cannot read data from workspace file reason: %w", err)
	}

	return nil
}

// AbsPath returns you absolute path of given ref. e.g. path to '//my/module'.
func (ws *Workspace) AbsPath(r string) string {
	return Ref(r).AbsPath(ws.rootDir)
}
