package yaml

import (
	"github.com/unravela/artisan/api"
	"testing"
)

func TestLoadModule(t *testing.T) {
	// when we lood the test MODULE.yaml
	var module api.Module
	err := LoadModule("./testdata/MODULE.yaml", &module)
	if err != nil {
		t.FailNow()
	}

	// then the module should have 2 tasks
	if len(module.Tasks) != 2 {
		t.FailNow()
	}

	// ... and 'build' task is present with './build.sh' script
	buildTask := module.Task("build")
	if buildTask == nil {
		t.FailNow()
	}

	if buildTask.Script != "./build.sh" {
		t.FailNow()
	}
}
