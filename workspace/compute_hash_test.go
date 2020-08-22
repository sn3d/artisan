package workspace

import (
	"fmt"
	"os"
	"testing"
)

func Test_Included(t *testing.T) {
	includes := []string {
		"/src",
	}

	res := included("hello/world", includes);
	if res {
		t.Errorf("expected res is 'false'")
	}

	res = included("hello/src", includes)
	if res {
		t.Errorf("expected res is 'false'")
	}

	res = included("/src", includes); fmt.Println(res)
	if !res {
		t.Errorf("expected res is 'true'")
	}

	res = included("/src/some/file.txt", includes); fmt.Println(res)
	if !res {
		t.Errorf("expected res is 'true'")
	}
}

func Test_IncludedAll(t *testing.T) {
	includes := []string { "" }
	res := included("hello/world", includes);
	if !res {
		t.Errorf("expected res is 'true'")
	}
}

func Test_Walk(t *testing.T) {
	// given includes and exludes for frontend
	includes := []string {
		"frontend",
	}

	excludes := []string {
		"frontend/dist",
		"frontend/node_modules",
	}

	// when we walk through 'go_with_javascript' workspace
	count := 0
	walk("testdata/compute_hash_test", "", includes, excludes, func(path string, f os.FileInfo) {
		count++
		fmt.Println(path)
	})

	// when we should have only files from 'frontend' but no 'dist' or 'node_modules'
	if count != 2 {
		t.Errorf("Somwthing wrong during walk")
	}
}