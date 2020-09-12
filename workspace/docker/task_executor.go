package docker

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
)

// TaskExecutor executes tasks via Docker
type TaskExecutor struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

// NewEnvTaskExecutor initialize faction builder with docker client based on
// docker env. variables like DOCKER_HOST etc..
func NewEnvTaskExecutor() (*TaskExecutor, error) {
	var executor *TaskExecutor
	var err error

	if dockr, err := client.NewEnvClient(); err == nil {
		executor = &TaskExecutor{
			Docker: dockr,
		}
	}

	return executor, err
}

// Exec executes the given task in given image
func (e *TaskExecutor) Exec(t *api.Task, img *api.Image, buildDir string) error {

	fmt.Printf(" - %s: Preparing...", t.Ref)
	ctx := context.Background()
	workingDir := "/build/" + t.Ref.GetPath()
	resp, err := e.Docker.ContainerCreate(
		ctx,
		&container.Config{
			Image:      img.ID,
			Cmd:        t.Cmd,
			WorkingDir: workingDir,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: buildDir,
					Target: "/build",
				},
			},
		},
		&network.NetworkingConfig{},
		"")

	if err != nil {
		fmt.Printf("\033[2K\r - %s: [ERROR]\n", t.Ref)
		return err
	}

	fmt.Printf("\033[2K\r - %s: Executing...", t.Ref)
	err = e.Docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		fmt.Printf("\033[2K\r - %s: [ERROR]\n", t.Ref)
		return err
	}

	out, err := e.Docker.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Printf("\033[2K\r - %s: %s", t.Ref, scanner.Text())
	}

	fmt.Printf("\033[2K\r - %s: [DONE]\n", t.Ref)

	return nil
}
