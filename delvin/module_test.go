package delvin

import "testing"

func TestLoadModule(t *testing.T) {
	// when we load valid HCL module file
	module := &Module{
		Ref: "//apps/webapp",
	}
	err := loadModuleHcl("../testdata/simplerepo", module)

	// then no error is occurred
	if err != nil {
		t.Errorf("We've got error %v", err)
	}

	// ... and the module must have tasks loaded
	if len(module.Tasks) != 2 {
		t.Errorf("We're expecting 2 modules but we've got %d", len(module.Tasks))
	}

	// ... and one task is "jdk8" class named as "build"
	buildTask := module.Task("build")
	if buildTask == nil {
		t.Errorf("there is no 'build' task!")
	}

	if buildTask.Class != "jdk8" {
		t.Errorf("the build task is not 'jdk8' class")
	}
}
