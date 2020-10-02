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
	jdk11Fact := ws.Factions["jdk11"]
	if jdk11Fact == nil {
		t.FailNow()
	}

	if jdk11Fact.Src != "//envs/jdk11" {
		t.FailNow()
	}

	// ... and workspace must have 'python3' env. with image defined
	pythonFact := ws.Factions["python3"]
	if pythonFact == nil {
		t.FailNow()
	}

	if pythonFact.Image != "python:3.7-alpine" {
		t.FailNow()
	}
}
