package workspace

import (
	"encoding/binary"
	"github.com/unravela/artisan/api"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func computeTaskHash(ws *Workspace, task *api.Task) api.TaskHash {
	buf := fnv.New64a()

	includes := task.Include
	if len(includes) == 0 {
		includes = []string{""}
	}

	excludes := task.Exclude
	if len(excludes) == 0 {
		excludes = []string{}
	}

	taskDir := ws.AbsPath(task.Ref)
	walk(taskDir, "", includes, excludes, func(path string, fi os.FileInfo) {
		buf.Write([]byte(path))

		// convert time to binary data and write
		t := fi.ModTime().Unix()
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(t))
		buf.Write(b)
	})

	hash := buf.Sum64()
	return api.TaskHash(hash)
}

func walk(taskDir string, subDir string, includes []string, excludes []string, walkFun func(string, os.FileInfo)) error {
	dir := filepath.Join(taskDir, subDir)
	content, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range content {
		path := filepath.ToSlash(filepath.Join(subDir, f.Name()))

		isIncluded := included(path, includes)
		isExcluded := excluded(path, excludes)

		if isIncluded && !isExcluded {
			if !f.IsDir() {
				// found file
				walkFun(path, f)
			} else {
				if err := walk(taskDir, path, includes, excludes, walkFun); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func included(path string, includes []string) bool {
	for _, incl := range includes {
		if strings.HasPrefix(path, incl) {
			return true
		}
		if strings.HasPrefix(incl, path) {
			return true
		}
	}
	return false
}

func excluded(path string, excludes []string) bool {
	for _, exl := range excludes {
		if strings.HasPrefix(path, exl) {
			return true
		}
	}
	return false
}
