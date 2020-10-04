package artisan

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)


// This function traverse upside over parent directories of given path until we
// found the folder that is matching to 'match' function
func findRoot(path string, match func(string) bool) (string, error) {
	var (
		currentDir string
		err error
	)

	currentDir, err = filepath.Abs(path)
	for {
		// if directory match the root criteria e.g. contain
		// MODULE.* file or WORKSPACE.* file.
		if match(currentDir) {
			return currentDir, nil
		}

		prevDir := currentDir
		currentDir, err = filepath.Abs(currentDir + "/..")

		// this check is because when we're on top (e.g. '/' or 'c:/')
		// we need to break the loop. Without this break, the find will enter
		// into infinity loop where currentDir is always root dir.
		if prevDir == currentDir {
			break
		}

		// ... or some unknown error is raised
		if err != nil {
			break
		}
	}

	return "", errors.New("Root not found")
}

// isModuleDir is used as match parameter for findRoot and check if directory
// contains one of the MODULE.* files (e.g. MODULE.hcl or MODULE.yaml, doesn't
// matter on extension).
func isModuleDir(dir string) bool {
	return hasFilePrefix(dir, "MODULE.")
}

// isWorkspaceDir is used as match parameteer for findRoot and check if
// directory contains one of the WORKSPACE.* files (e.g. WORSPACE.hcl or
// WORKSPACE.yaml, doesn't matter on extension).
func isWorkspaceDir(dir string) bool {
	return hasFilePrefix(dir, "WORKSPACE.")
}

func hasFilePrefix(dir string, prefix string) bool {
	f,err := os.Open(dir)
	if err != nil {
		return false
	}
	defer f.Close()

	files, err := f.Readdir(0)
	if err != nil {
		return false
	}

	for i := range files {
		if !files[i].IsDir() {
			if strings.HasPrefix(files[i].Name(), prefix) {
				return true
			}
		}

	}

	return false
}