package docker

import (
	"github.com/docker/docker/client"
	"github.com/unravela/delvin/delvin"
)

type TaskExecutor struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

func (e *TaskExecutor) Exec(t *delvin.Task) {

}
