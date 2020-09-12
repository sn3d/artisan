package api

// TaskExecutor provide execution implementation logic
type TaskExecutor interface {
	Exec(t *Task, img *Image, buildDir string) error
}
