package api

// EnvironmentDef describe the environment in which will be task executed.
//
// Let's imagine you want to build NodeJS project. You can define 'build' task
// that is associated to 'nodejs' environment. This env. is basically the
// docker image where project will be mounted and build.
//
// Example of 'nodejs' env. in artisan file:
//
//    environment "nodejs" {
//       src: "//envs/nodejs
//    }
//
// All 'nodejs' participants will be build within docker image that is
// defined in /envs/nodejs/Dockerfile.
//
// Now, you can create build task associated with this environment
//
//    task "nodejs" "build" {
//       ...
//    }
//
// Environment names they're not defined in WORKSPACE file are resolved as
// standard docker images. Following example means the 'build' task will be
// executed in directly in 'golang:1.15.0-buster' docker
// image
//
// Example:
//
//     task "golang:1.15.0-buster" "build" {
//         ...
//     }
type EnvironmentDef struct {
	// Name of the environment
	Name string `hcl:"name,label"`
	// Src is ref. to directory where is source of Environment (e.g. Dockerfile)
	Src string `hcl:"src,optional"`
	// Image is used when we want to use image from official docker hub
	Image string `hcl:"image,optional"`
}

// Environments is set of env items. Because we don't want duplicities
// for name, it's implemented as "set" by name.
type EnvironmentDefs map[string]*EnvironmentDef

// EnvironmentID refers to environment resolved from EnvironmentDef. For task
// execution we need to have environment (docker image), not just definition.
// The resolving is happening in EnvironmentRegistry
type EnvironmentID string

// EnvironmentRegistry resolve env. definition to real environment identified
// by EnvironmentID
type EnvironmentRegistry interface {
	Build(envDef *EnvironmentDef, srcDir string) (EnvironmentID, error)
}

func (e *EnvironmentDef) String() string {
	return e.Name
}

// NewEnvironments convert array to set of Environments
func NewEnvironments(arr []*EnvironmentDef) EnvironmentDefs {
	envs := make(EnvironmentDefs)
	for _, c := range arr {
		envs[c.Name] = c
	}
	return envs
}

