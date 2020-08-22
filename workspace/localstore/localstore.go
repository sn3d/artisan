package localstore

import (
	"github.com/unravela/delvin/api"
	"os"
)

type LocalStore struct {
	// path to local storage directory where are files for task hashes etc..
	dir string

	// this map hold task hashes where key is hash of the unique hash of ref. and value
	// is hash of all files they're belong to task
	hashes map[uint64]api.TaskHash
}

// Open the localstore
func Open(dir string) (*LocalStore, error) {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, os.ModeDir)
	}

	hashes, err := loadTaskHashes(dir)
	if err != nil {
		return nil, err
	}

	store := &LocalStore{
		dir:    dir,
		hashes: hashes,
	}

	return store, nil
}
