package yaml

import (
	"github.com/unravela/artisan/api"
	"testing"
)

func TestLoadWorkspace(t *testing.T) {
	// when we lood the test WORKSPACE.yaml
	var ws api.Workspace
	err := LoadWorkspace("./testdata/WORKSPACE.yaml", &ws)
	if err != nil {
		t.FailNow()
	}

	// then the workspace must have 'jdk11' environment with src defined
	jdk11Env := ws.Environments["jdk11"]
	if jdk11Env == nil {
		t.FailNow()
	}

	if jdk11Env.Src != "//envs/jdk11" {
		t.FailNow()
	}

	// ... and workspace must have 'python3' env. with image defined
	pythonEnv := ws.Environments["python3"]
	if pythonEnv == nil {
		t.FailNow()
	}

	if pythonEnv.Image != "python:3.7-alpine" {
		t.FailNow()
	}
}
