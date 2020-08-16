package delvin

import (
	"strings"
	"testing"
)

func TestRef_GetType(t *testing.T) {
	// given ref. with 'git' type
	ref := Ref("remote://app/custom:build")

	// when we get type
	typ := ref.GetType()

	// then type must be 'git'
	if typ != "remote" {
		t.FailNow()
	}
}

func TestRef_GetPath(t *testing.T) {
	// given complicated ref.
	ref := Ref("remote://app/custom:build")

	// when we get type
	path := ref.GetPath()

	// then path must be 'app/custom'
	if path != "app/custom" {
		t.FailNow()
	}
}

func TestRef_SetTask(t *testing.T) {
	// when we append task
	ref := Ref("//apps/webapp").SetTask("build")
	if ref != "//apps/webapp:build" {
		t.Errorf("Task 'build' wasn't appended correctly")
	}

	// when we change the existing task
	ref = Ref("//apps/webapp:build").SetTask("test")
	if ref != "//apps/webapp:test" {
		t.Errorf("Task wasn't changed to 'test'")
	}
}

func TestRef_GetTask(t *testing.T) {
	// when we have simple ref '//my/app:build'
	task := Ref("//my/app:build").GetTask()

	// then we should get the 'build'
	if task != "build" {
		t.FailNow()
	}
}

func TestRef_GetTask2(t *testing.T) {
	// when we have just task ':build'
	task := Ref(":build").GetTask()

	// then we should get the 'build'
	if task != "build" {
		t.FailNow()
	}
}

func TestRef_GetTask3(t *testing.T) {
	// when we have full ref 'type://app:build'
	task := Ref("type://app:build").GetTask()

	// then we should get the 'build'
	if task != "build" {
		t.FailNow()
	}
}

func TestRef_GetTask4(t *testing.T) {
	// when we have ref without task 'type://app'
	task := Ref("type://app").GetTask()

	// then task should be empty string
	if task != "" {
		t.FailNow()
	}
}

func TestRef_AbsPath(t *testing.T) {
	// when we get absolute path for ref. //my/module
	path := Ref("//my/module").AbsPath("../testdata")

	// then the path must point to testdata/my/module
	if !strings.HasSuffix(path, "/testdata/my/module") {
		t.Errorf("We expect the path has suffix '/testdata/my/module'. Check the path %s", path)
	}
}
