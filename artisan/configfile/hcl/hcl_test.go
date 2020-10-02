package hcl

import (
	"github.com/unravela/artisan/api"
	"testing"
)

func TestLoadModule(t *testing.T) {
	// when we load valid HCL module file
	module := &api.Module{
		Ref: "//apps/webapp",
	}
	err := LoadModule("testdata/MODULE.hcl", module)

	// then no error is occurred
	if err != nil {
		t.Errorf("We've got error %v", err)
	}

	// ... and one task is "go" faction named as "build"
	buildTask := module.Task("build")
	if buildTask == nil {
		t.Errorf("there is no 'build' task!")
	}

	if buildTask.FactionName != "go-1.13" {
		t.Errorf("the build task is not 'jdk8' faction")
	}
}
