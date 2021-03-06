package api

import (
	"hash/fnv"
	"path/filepath"
	"strings"
)

// Ref is reference to module or task. The format looks like 'artisan://path/to/module:task'.
// If artisan or path is not defined, we considering the 'current'. That means current artisan,
// or current module. For instance ':build' in task dependencies refers to same module the dependant
// task is.
type Ref string


// StringToRef is used for converting any string to Ref.
// Don't use direct casting even Ref is string and prefer this function
// because some pre-processing and validation
func StringToRef(str string) Ref {
	trimmed := strings.Trim(str, " \t")
	return Ref(trimmed)
}

// NewRef construct the reference with artisan, path and task.
// If ref. cannot be constructed because e.g. path is empty string,
// the output will be empty Ref
func NewRef(ws string, path string, task string) Ref {
	refStr := path

	if ws != "" {
		refStr = ws + ":" + path
	}

	if task != "" {
		refStr = refStr + ":" + task
	}

	return Ref(refStr)
}

// GetWorkspace parse the Ref and returns you 'artisan' part.
// If it's empty string, that means 'this' artisan
func (r Ref) GetWorkspace() string {
	var typ = ""
	idx := strings.Index(string(r), "://")
	if idx > 0 {
		typ = string(r)[0:idx]
	}
	return typ
}

// GetPath returns you path part in ref which is 'path/to/subdir' in 'ws://path/to/subdir:task'
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
// placed as last part in ref 'ws://path/to/app:task'
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
	if r.GetWorkspace() != "" {
		out += r.GetWorkspace() + ":"
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

// AbsPath function gives you absolute path of
// reference in form '/abspath/to/module` on Windows and Unix systems.
// This function need artisan root directory for determining absolute path.
func (r Ref) AbsPath(rootDir string) string {
	path := filepath.Join(rootDir, r.GetPath())
	path, _ = filepath.Abs(path)
	path = filepath.ToSlash(path)
	return path
}

// IsOnlyTask returns you true if Ref is in form ":build"
func (r Ref) IsOnlyTask() bool {
	return strings.HasPrefix(string(r), ":")
}
