package api

// Environment determine in which environment will be task performed.
//
// Let's imagine you want to build NodeJS project. You can define 'build' task
// that is associated to 'nodejs' environment. This env. is basically definition of
// docker image where project will be mounted and build.
//
// Example of 'nodejs' env. in artisan file:
//
//    environment "nodejs" {
//       src: "//envs/nodejs
//    }
//
// All 'nodejs' participants will be build within docker image that is defined in /envs/nodejs/Dockerfile.
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
//     task "golang:1.15.0-buster" "build" {
//         ...
//     }
type Environment struct {
	// Name of the environment
	Name string `hcl:"name,label"`
	// Src is ref. to directory where is source of Environment (e.g. Dockerfile)
	Src string `hcl:"src,optional"`
	// Image is used when we want to use image from official docker hub
	Image string `hcl:"image,optional"`
}

// Environments is set of env items. Because we don't want duplicities
// for name, it's implemented as "set" by name.
type Environments map[string]*Environment

func (e *Environment) String() string {
	return e.Name
}

// NewEnvironments convert array to set of Environments
func NewEnvironments(arr []*Environment) Environments {
	envs := make(Environments)
	for _, c := range arr {
		envs[c.Name] = c
	}
	return envs
}
