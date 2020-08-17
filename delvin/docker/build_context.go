package docker

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type buildContextOptions struct {
	root     string
	includes []string
}

// creates build context in-memory and returns you reader. The build context is
// an TAR archive which contains  Dockerfile and all files needed for creating docker image.
//
func createBuildContext(opts buildContextOptions) io.Reader {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	filepath.Walk(opts.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() {
			slashPath := filepath.ToSlash(path)[len(opts.root):]
			if matchToIncludes(opts.includes, slashPath) {
				appendFile(opts.root, path, tw)
			}
		}

		return nil
	})

	tw.Close()
	return &buf
}

func matchToIncludes(includes []string, p string) bool {
	// if includes is not set, we will include everything
	if len(includes) == 0 {
		return true
	}

	for _, incl := range includes {
		matched, _ := path.Match(incl, p)
		if matched {
			return true
		}
	}
	return false
}

func appendDir(root string, dir string, tw *tar.Writer) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if !file.IsDir() {
			appendFile(root, filepath.Join(dir, file.Name()), tw)
		} else {
			appendDir(root, filepath.Join(dir, file.Name()), tw)
		}
	}
}

func appendFile(root string, filename string, tw *tar.Writer) {

	data, _ := ioutil.ReadFile(filename)
	fname := filepath.ToSlash(filename[len(root):])
	tw.WriteHeader(&tar.Header{
		Name: fname,
		Mode: 0750,
		Size: int64(len(data)),
	})
	tw.Write(data)
}
