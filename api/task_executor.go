package api

type TaskExecutor interface {
	Exec(t *Task, img *Image, buildDir string) error
}
