package localstore

import (
	"os"
	"path/filepath"

	"github.com/unravela/artisan/api"
)

const (
	// DirName is name of the directory used by localstore. This directory
	// is usually placed in root of your workspace.
	DirName = ".artisan"
)

// LocalStore instance provide you access to all working data like task hashes etc..
type LocalStore struct {
	// path to local storage directory where are files for task hashes etc..
	dir string

	// this map hold task hashes where key is hash of the unique hash of ref. and value
	// is hash of all files they're belong to task
	hashes map[uint64]api.TaskHash
}

// Open the localstore which is part of the given
// workspace root directory. The 'rootDir' is path
// which not include the localstore directory.
func Open(rootDir string) (*LocalStore, error) {

	dir := filepath.Join(rootDir, DirName)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.Mkdir(dir, os.ModeDir|0750)
		if err != nil {
			return nil, err
		}
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
