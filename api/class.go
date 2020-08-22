package api

// Class or builder class, determine in which build container the task will be performed. Let's imagine you
// have task that is building your NodeJS application. We need specific version of NodeJS. The task is
// attached to class 'nodejs' which is docker image with specific version of NodeJS:
// for this task could be 'node:12.13.0' docker image.
//
//    class "nodejs" {
//        image: "nodejs:12.13.0"
//    }
//
// Also you can define your own Dockerfile for class. You will define path to this docker file.
//
//    class "gradle" {
//        src: "//forges/gradle/Dockerfile"
//    }
//
type Class struct {
	// Name of the class
	Name string `hcl:"name,label"`
	// Src is ref. to directory where is source of Forge (e.g. Dockerfile)
	Src string `hcl:"src,optional"`
	// Image is used when we want to use Forge from docker hub
	Image string `hcl:"image,optional"`
}

func (c *Class) String() string {
	return c.Name
}
