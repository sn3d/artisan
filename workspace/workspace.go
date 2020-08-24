package workspace

import (
	"errors"
	"fmt"
	"github.com/unravela/delvin/api"
	"github.com/unravela/delvin/workspace/docker"
	"github.com/unravela/delvin/workspace/localstore"
	"os"
	"path/filepath"
	"strings"
)

type Workspace struct {
	// absolute path to root of the workspace. This is a dir where
	// workspace file is situated
	rootDir string

	// hold all available factions for workspace
	factions api.Factions

	// mainModule is module with reference '//'
	mainModule *api.Module
}

// Open opens the workspace. If you're in subdirectory of workspace, the
// function will traverse upside until workspace root is found.
func Open(currentDir string) (*Workspace, error) {
	rootDir, err := findWorkspaceRoot(currentDir)
	if err != nil {
		return nil, errors.New("Cannot open workspace")
	}

	ws := &Workspace{
		rootDir: rootDir,
	}

	hclFile := filepath.Join(rootDir, WorkspaceFile)
	if err = loadWorkspaceFromHCL(hclFile, ws); err != nil {
		return nil, err
	}

	return ws, nil
}

// Faction resolve a faction def. for the given name
func (ws *Workspace) Faction(name string) *api.Faction {
	for _, c := range ws.factions {
		if c.Name == name {
			return c
		}
	}

	// faction not found - let's transform it into Docker image
	defaultFaction := &api.Faction{
		Name:  name,
		Image: name,
	}

	return defaultFaction
}

// Module returns a module for given ref, no matter if
// ref. contains also task. If the module is root '//', then
// it's returned main module within workspace file
func (ws *Workspace) Module(ref api.Ref) (*api.Module, error) {
	var module *api.Module

	if ref.GetPath() == "" {
		module = ws.mainModule
	} else {
		module = &api.Module{
			Ref: ref,
		}

		hclFile := filepath.Join(ws.AbsPath(ref), ModuleFile)
		err := LoadModuleFromHCL(hclFile, module)
		if err != nil {
			return nil, err
		}
	}

	return module, nil
}

// FindModule find a module for given path. This
// function will traverse up by parent directories until the module
// is found.
func (ws *Workspace) FindModule(path string) api.Ref {
	root, err := findRootFor(path, ModuleFile)
	if err != nil {
		return ""
	}
	return ws.AbsPathToRef(root)
}

// Task gives you task by given reference. If there is no
// task, the null is returned.
func (ws *Workspace) Task(ref api.Ref) (*api.Task, error) {
	tskName := ref.GetTask()
	if tskName == "" {
		return nil, fmt.Errorf("Task is not defined")
	}

	module, err := ws.Module(ref)
	if err != nil {
		return nil, fmt.Errorf("Module for task not fount")
	}

	task := module.Task(tskName)
	if task == nil {
		return nil, fmt.Errorf("Task not present in module")
	}

	return task, nil
}

// this function traverse upside over directories until we found directory with workspace file.
// The reason why we need to traverse-up the directories is because
// user might call command in any subdirectory of repository.
//
// Function is returning absolute path of root directory or empty string with
// error if there is no workspace file.
func findWorkspaceRoot(path string) (string, error) {
	return findRootFor(path, WorkspaceFile)
}

// This function traverse upside over parent directories of given path until we
// found the folder with the given filename inside.
func findRootFor(path string, filename string) (string, error) {
	currentDir, _ := filepath.Abs(path)

	for {
		wsFile := filepath.Join(currentDir, filename)
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
			// we need to break the loop. Without this break, the find will enter
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

	return "", errors.New("Root of woskpace not found")
}

// AbsPath returns you absolute path of given ref. e.g. path to '//my/module'.
func (ws *Workspace) AbsPath(r api.Ref) string {
	path := filepath.Join(ws.rootDir, r.GetPath())
	path, _ = filepath.Abs(path)
	path = filepath.ToSlash(path)
	return path
}

// AbsPathToRef todo
func (ws *Workspace) AbsPathToRef(abspath string) api.Ref {
	if !strings.HasPrefix(abspath, ws.rootDir) {
		return ""
	}

	refPath := abspath[len(ws.rootDir):]
	refPath = "/" + filepath.ToSlash(refPath)
	return api.Ref(refPath)
}

// Run perform the given task and task's dependencies.
func (ws *Workspace) Run(taskRef api.Ref) error {
	var engine api.Engine
	var task *api.Task
	var err error

	lstore, err := localstore.Open(filepath.Join(ws.rootDir, ".delvin"))
	if err != nil {
		return err
	}

	if err = docker.SetupEngine(&engine); err != nil {
		return err
	}

	if task, err = ws.Task(taskRef); err != nil {
		return err
	}

	// phase 1: ensure images

	fmt.Println("(1/2) Resolve images")

	allTasks := topoSort(task, ws)
	allFactions := allTasks.GetFactions(ws)
	allImages := make(Images)
	for _, fact := range allFactions {
		imageSrcDir := ws.AbsPath(api.Ref(fact.Src))
		img, err := engine.Registry.Build(fact, imageSrcDir)
		if err != nil {
			return err
		}
		allImages[fact.Name] = img
	}

	// phase 2: run tasks

	fmt.Println("(2/2) Run tasks")
	for _, tsk := range allTasks {
		// check if task neet to be build
		if isUpToDate(task, lstore, ws) {
			fmt.Printf(" - %s: [UP TO DATE]\n", tsk.Ref)
			break
		}

		img := allImages[tsk.FactionName]
		err := engine.Executor.Exec(tsk, img, ws.rootDir)
		if err != nil {
			return err
		}

		// store the new hash
		newHash := computeTaskHash(ws, tsk)
		lstore.PutTaskHash(tsk.Ref, newHash)
	}

	return nil
}

func isUpToDate(tsk *api.Task, lstore *localstore.LocalStore, ws *Workspace) bool {
	storedHash := lstore.GetTaskHash(tsk.Ref)
	actualHash := computeTaskHash(ws, tsk)
	if storedHash == actualHash {
		return true
	}
	return false
}
