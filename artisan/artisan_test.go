package artisan_test

import (
	"github.com/unravela/artisan/artisan"
	"path/filepath"
	"testing"
)

func TestOpen(t *testing.T) {
	// when we try to open artisan from subdirectory
	ws, err := artisan.OpenWorkspace("testdata/workspace_test/apps/frontend")

	// then no error is occurred
	if err != nil {
		t.Errorf("Some error occured when we wanted open the artisan")
	}

	// ... and artisan need to be initialized
	if ws == nil {
		t.Errorf("Some error occured when opening the artisan ")
	}

	// ... and we can reach the faction
	fact := ws.Faction("@go")
	if fact == nil {
		t.Errorf("We're expecting '@go' fact available.")
	}
}

func TestOpenInvalidFolder(t *testing.T) {
	// when we try to open random folder out of the artisan
	ws, err := artisan.OpenWorkspace("testdata")

	// then error is occurred
	if err == nil {
		t.Errorf("Weird, there is no error")
	}

	// ... and no artisan is open
	if ws != nil {
		t.Errorf("Some artisan was open")
	}
}

func TestWorkspace_Task(t *testing.T) {
	// given artisan of simple repository
	ws, err := artisan.OpenWorkspace("testdata/workspace_test")
	if err != nil {
		t.Errorf("Cannot open repository")
	}

	// when we get the existing task
	tsk, _ := ws.Task("//apps/backend:build")
	if tsk == nil {
		t.Errorf("Cannot find task!")
	}

	// then the task must be returned with correct reference.
	if tsk.Ref != "//apps/backend:build" {
		t.Errorf("The ref. of the task doesn't match")
	}
}

func TestWorkspace_FindModule(t *testing.T) {
	// given open artisan
	ws, _ := artisan.OpenWorkspace("testdata/workspace_test")

	// when we look for module of /apps/frontend/src folder
	mref := ws.FindModule("testdata/workspace_test/apps/frontend/src")

	// then the module is //apps/frontend
	if mref != "//apps/frontend" {
		t.Errorf("Cannot get the module")
	}
}

func TestWorkspace_AbsPathToRef(t *testing.T) {
	// when we get the ref for absolute path ./apps/backend
	ws, _ := artisan.OpenWorkspace("testdata/workspace_test")
	abspath, _ := filepath.Abs("testdata/workspace_test/apps/backend")
	ref := ws.AbsPathToRef(abspath)

	// the ref must be //apps/backend
	if ref != "//apps/backend" {
		t.Errorf("invalid ref.")
	}
}
