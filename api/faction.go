package api

// Faction determine in which environment, whit what tools will be task performed.
//
// Let's imagine you want to build NodeJS project. You can define 'build' task
// that is associated to '@nodejs' faction. This faction is basically definition of
// docker image where project will be mounted and build.
//
// Example of '@nodejs' faction in workspace file:
//
//    faction "@nodejs" {
//       src: "//factions/nodejs
//    }
//
// All 'nodejs' participants will be build within docker image that is defined in /factions/nodejs/Dockerfile.
// Now, you can create build task associated with this faction
//
//    task "@nodejs" "build" {
//       ...
//    }
//
// All local factions start with at '@' character. Faction names started without
// this character are resolved as standard docker images. Following example means
// the 'build' task will be executed in directly in 'golang:1.15.0-buster' docker
// image
//
// Example:
//     task "golang:1.15.0-buster" "build" {
//         ...
//     }
type Faction struct {
	// Name of the faction
	Name string `hcl:"name,label"`
	// Src is ref. to directory where is source of Faction (e.g. Dockerfile)
	Src string `hcl:"src,optional"`
	// Image is used when we want to use Forge from docker hub
	Image string `hcl:"image,optional"`
}

// Factions is set of faction items. Because we don't want duplicities
// for name, it's implemented as "set" by name.
type Factions map[string]*Faction

func (f *Faction) String() string {
	return f.Name
}

// NewFactions convert array to set of 'Factions'
func NewFactions(arr []*Faction) Factions {
	factions := make(Factions)
	for _, c := range arr {
		factions[c.Name] = c
	}
	return factions
}


