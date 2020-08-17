package delvin_test

import (
	"github.com/unravela/delvin/delvin"
	"testing"
)

func TestOpen(t *testing.T) {
	// when we try to open workspace from subdirectory
	ws, err := delvin.Open("../testdata/simplerepo/webapp")

	// then no error is occurred
	if err != nil {
		t.Errorf("Some error occured when we wanted open the workspace")
	}

	// ... and workspace need to be initialized
	if ws == nil {
		t.Errorf("Some error occured when opening the workspace ")
	}

	// ... and we can reach the forge from file
	class := ws.Class("nodejs")
	if class == nil {
		t.Errorf("We're expecting 'nodejs' builder class available.")
	}
}

func TestOpenInvalidFolder(t *testing.T) {
	// when we try to open random folder out of the workspace
	ws, err := delvin.Open("../testdata")

	// then error is occurred
	if err == nil {
		t.Errorf("Weird, there is no error")
	}

	// ... and no workspace is open
	if ws != nil {
		t.Errorf("Some workspace was open")
	}
}

func TestWorkspace_Task(t *testing.T) {
	// given workspace of simple repository
	ws, err := delvin.Open("../testdata/simplerepo")
	if err != nil {
		t.Errorf("Cannot open repository")
	}

	// when we get the existing task
	tsk := ws.Task("//apps/webapp:test")
	if tsk == nil {
		t.Errorf("Cannot find task!")
	}

	// then the task must be returned with correct reference.
	if tsk.Ref != "//apps/webapp:test" {
		t.Errorf("The ref. of the task doesn't match")
	}
}
