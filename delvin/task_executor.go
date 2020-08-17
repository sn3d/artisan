package delvin

type TaskExecutor interface {
	Exec(t *Task, img *ClassImage, buildDir string) error
}
