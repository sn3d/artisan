package delvin

import (
	"hash/fnv"
	"path/filepath"
	"strings"
)

// Ref is reference to module or task. The format is borrowed from
// Bazel and looks like 'type://path/to/subdir:task', where type and task is optional.
// Normally type is not set and the ref looks like something you know from bazel '//app:task'.
type Ref string

// GetType parse the Ref and returns you 'type' part.
func (r Ref) GetType() string {
	var typ = ""
	idx := strings.Index(string(r), "://")
	if idx > 0 {
		typ = string(r)[0:idx]
	}
	return typ
}

// GetPath returns you path part in ref which is 'path/to/subdir' in 'type://path/to/subdir:task'
func (r Ref) GetPath() string {
	var path = ""
	start := strings.Index(string(r), "//")
	if start >= 0 {
		path = string(r)[(start + len("//")):]
		end := strings.Index(path, ":")
		if end >= 0 {
			path = path[:end]
		}
	}

	return path
}

// GetTask returns you task, if there is some task present. The task is
// placed as last part in ref 'type://path/to/app:task'
func (r Ref) GetTask() string {
	start := strings.LastIndex(string(r), ":")
	task := string(r)[start+1:]

	if strings.HasPrefix(task, "//") {
		// there is no task because task cannot start with '//'
		return ""
	}
	return task
}

// SetTask append task to the end of the reference and return new ref. value.
// The old task will be changed.
func (r Ref) SetTask(task string) Ref {
	out := ""
	if r.GetType() != "" {
		out += r.GetType() + ":"
	}

	out += "//" + r.GetPath() + ":" + task
	return Ref(out)
}

func (r Ref) String() string {
	return string(r)
}

// GetHash returns hash of ref.
func (r Ref) GetHash() uint64 {
	h := fnv.New64a()
	h.Write([]byte(r))
	return h.Sum64()
}

// this function gives you absolute path of
// reference in form '/abspath/to/module` on Windows and Unix systems.
// This function need workspace root directory for determining absolute path.
func (r Ref) AbsPath(rootDir string) string {
	path := filepath.Join(rootDir, r.GetPath())
	path, _ = filepath.Abs(path)
	path = filepath.ToSlash(path)
	return path
}
