package localstore_test

import (
	"github.com/unravela/artisan/workspace/localstore"
	"io/ioutil"
	"testing"
)

func TestLocalStore_GetAndPut(t *testing.T) {
	// given open local storage
	dir, _ := ioutil.TempDir("", "delvin")

	// when we put some task hash
	ls, _ := localstore.Open(dir)
	if ls == nil {
		t.Errorf("Cannot open local storage")
	}
	err := ls.PutTaskHash("//frontend:test", 123)
	if err != nil {
		t.Errorf("Cannot put the task")
	}

	// then we should be able get the hash back
	ls2, _ := localstore.Open(dir)
	if ls2 == nil {
		t.Errorf("Cannot open local storage")
	}

	hash := ls2.GetTaskHash("//frontend:test")
	if hash != 123 {
		t.Errorf("Error getting hash")
	}

}
