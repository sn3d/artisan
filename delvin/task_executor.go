package delvin

type TaskExecutor interface {
	Exec(t *Task)
}
