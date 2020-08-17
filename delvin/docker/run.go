package docker

import (
	"github.com/docker/docker/client"
	"github.com/unravela/delvin/delvin"
)

// SetupRunContext function setup the context with docker
// executor and class builder.
func SetupRunContext(ctx *delvin.RunContext) error {

	dockr, err := client.NewEnvClient();
	if err != nil {
		return err
	}

	ctx.Builder = &ClassBuilder{
		Docker: dockr,
	}

	ctx.Executor = &TaskExecutor{
		Docker: dockr,
	}

	return nil
}
