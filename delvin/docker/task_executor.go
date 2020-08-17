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
	"github.com/unravela/delvin/delvin"
)

type TaskExecutor struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

func (e *TaskExecutor) Exec(t *delvin.Task, img *delvin.ClassImage, buildDir string) error {
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
		return err
	}

	err = e.Docker.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}

	out, err := e.Docker.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     true,
	})

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	return nil
}
