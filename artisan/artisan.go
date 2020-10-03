package artisan

import (
	"errors"
	"fmt"
	"github.com/unravela/artisan/api"
	"github.com/unravela/artisan/artisan/configfile"
	"github.com/unravela/artisan/artisan/docker"
	"github.com/unravela/artisan/artisan/localstore"
	"os"
	"path/filepath"
	"strings"
)

// Artisan is main facade and provide most of the functionality over
// opened artisan
type Artisan struct {
	// absolute path to root of the workspace. This is a dir where
	// workspace file is situated
	//rootDir string

	// hold the opened workspace
	workspace api.Workspace
}

// Open opens the artisan. If you're in subdirectory of artisan, the
// function will traverse upside until artisan root is found.
func OpenWorkspace(currentDir string) (*Artisan, error) {
	rootDir, err := findWorkspaceRoot(currentDir)
	if err != nil {
		return nil, errors.New("Cannot open artisan")
	}

	instance := &Artisan{
		workspace: api.Workspace{
			RootDir: rootDir,
		},
	}

	if err = configfile.LoadWorkspace(rootDir, &instance.workspace); err != nil {
		return nil, err
	}

	return instance, nil
}

// Environment resolve a env. definition for the given name
func (inst *Artisan) Environment(name string) *api.Environment {
	for _, c := range inst.workspace.Environments {
		if c.Name == name {
			return c
		}
	}

	defaultEnvironment := &api.Environment{
		Name:  name,
		Image: name,
	}

	return defaultEnvironment
}

// Module returns a module for given ref, no matter if
// ref. contains also task. If the module is root '//', then
// it's returned main module within artisan file
func (inst *Artisan) Module(ref api.Ref) (*api.Module, error) {
	var module *api.Module

	if ref.GetPath() == "" {
		module = inst.workspace.MainModule
	} else {
		module = &api.Module{
			Ref: ref,
		}

		moduleDir := inst.AbsPath(ref)
		err := configfile.LoadModule(moduleDir, module)
		if err != nil {
			return nil, err
		}
	}

	return module, nil
}

// FindModule find a module for given path. This
// function will traverse up by parent directories until the module
// is found.
func (inst *Artisan) FindModule(path string) api.Ref {
	root, err := findRootFor(path, configfile.ModuleFile)
	if err != nil {
		return ""
	}
	return inst.AbsPathToRef(root)
}

// Task gives you task by given reference. If there is no
// task, the null is returned.
func (ws *Artisan) Task(ref api.Ref) (*api.Task, error) {
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

// this function traverse upside over directories until we found directory with artisan file.
// The reason why we need to traverse-up the directories is because
// user might call command in any subdirectory of repository.
//
// Function is returning absolute path of root directory or empty string with
// error if there is no artisan file.
func findWorkspaceRoot(path string) (string, error) {
	return findRootFor(path, configfile.WorkspaceFile)
}

// This function traverse upside over parent directories of given path until we
// found the folder with the given filename inside.
func findRootFor(path string, filename string) (string, error) {
	currentDir, _ := filepath.Abs(path)

	for {
		wsFile := filepath.Join(currentDir, filename)
		info, err := os.Stat(wsFile)

		// if we found artisan root
		if err == nil && info.Mode().IsRegular() {
			return currentDir, nil
		}

		// ... or there is no artisan file
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
func (inst *Artisan) AbsPath(r api.Ref) string {
	path := filepath.Join(inst.workspace.RootDir, r.GetPath())
	path, _ = filepath.Abs(path)
	path = filepath.ToSlash(path)
	return path
}

// AbsPathToRef todo
func (inst *Artisan) AbsPathToRef(abspath string) api.Ref {
	if !strings.HasPrefix(abspath, inst.workspace.RootDir) {
		return ""
	}

	refPath := abspath[len(inst.workspace.RootDir):]
	refPath = "/" + filepath.ToSlash(refPath)
	return api.Ref(refPath)
}

// Run perform the given task and task's dependencies.
func (inst *Artisan) Run(taskRef api.Ref) error {
	var engine api.Engine
	var task *api.Task
	var err error

	lstore, err := localstore.Open(inst.workspace.RootDir)
	if err != nil {
		return err
	}

	if err = docker.SetupEngine(&engine); err != nil {
		return err
	}

	if task, err = inst.Task(taskRef); err != nil {
		return err
	}

	// phase 1: ensure images

	fmt.Println("(1/2) Resolve images")

	allTasks := topoSort(task, inst)
	allEnvs := inst.extractEnvironments(allTasks)
	allImages := make(api.Images)
	for _, env := range allEnvs {
		imageSrcDir := inst.AbsPath(api.Ref(env.Src))
		img, err := engine.Registry.Build(env, imageSrcDir)
		if err != nil {
			return err
		}
		allImages[env.Name] = img
	}

	// phase 2: run tasks
	fmt.Println("(2/2) Run tasks")
	for _, tsk := range allTasks {
		// check if task neet to be build
		if isUpToDate(task, lstore, inst) {
			fmt.Printf(" - %s: [UP TO DATE]\n", tsk.Ref)
			break
		}

		img := allImages[tsk.EnvName]
		err := engine.Executor.Exec(tsk, img, inst.workspace.RootDir)
		if err != nil {
			return err
		}

		// store the new hash
		newHash := computeTaskHash(inst, tsk)
		lstore.PutTaskHash(tsk.Ref, newHash)
	}

	return nil
}

func isUpToDate(tsk *api.Task, lstore *localstore.LocalStore, inst *Artisan) bool {
	storedHash := lstore.GetTaskHash(tsk.Ref)
	actualHash := computeTaskHash(inst, tsk)
	if storedHash == actualHash {
		return true
	}
	return false
}

func (inst *Artisan) extractEnvironments(t api.Tasks) api.Environments {
	envs := make(api.Environments)
	for _, task := range t {
		envDef := inst.Environment(task.EnvName)
		envs[task.EnvName] = envDef
	}
	return envs
}

