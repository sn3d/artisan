package api

// TaskExecutor provide execution implementation logic
type TaskExecutor interface {
	Exec(t *Task, envID EnvironmentID, buildDir string) error
}
