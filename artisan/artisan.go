package artisan

import (
	"errors"
	"fmt"
	"github.com/unravela/artisan/api"
	"github.com/unravela/artisan/artisan/configfile"
	"github.com/unravela/artisan/artisan/docker"
	"github.com/unravela/artisan/artisan/localstore"
	"path/filepath"
	"strings"
)

// Artisan is main facade and provide most of the functionality over
// opened artisan
type Artisan struct {
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
func (inst *Artisan) Environment(name string) *api.EnvironmentDef {
	for _, c := range inst.workspace.Environments {
		if c.Name == name {
			return c
		}
	}

	defaultEnvironment := &api.EnvironmentDef{
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
	root, err := findRoot(path, isModuleDir)
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
	return findRoot(path, isWorkspaceDir)
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
	return api.StringToRef(refPath)
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
	allEnvIDs := make(map[string]api.EnvironmentID)
	for _, env := range allEnvs {
		imageSrcDir := inst.AbsPath(api.StringToRef(env.Src))
		envID, err := engine.Registry.Build(env, imageSrcDir)
		if err != nil {
			return err
		}
		allEnvIDs[env.Name] = envID
	}

	// phase 2: run tasks
	fmt.Println("(2/2) Run tasks")
	for _, tsk := range allTasks {
		// check if task neet to be build
		if isUpToDate(tsk, lstore, inst) {
			fmt.Printf(" - %s: [UP TO DATE]\n", tsk.Ref)
			continue
		}

		envID := allEnvIDs[tsk.EnvName]
		err := engine.Executor.Exec(tsk, envID, inst.workspace.RootDir)
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

func (inst *Artisan) extractEnvironments(t api.Tasks) api.EnvironmentDefs {
	envs := make(api.EnvironmentDefs)
	for _, task := range t {
		envDef := inst.Environment(task.EnvName)
		envs[task.EnvName] = envDef
	}
	return envs
}

