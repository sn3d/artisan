package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
)

// SetupEngine initialize engine with docker client
func SetupEngine(engine *api.Engine) error {

	dockr, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("cannot connect to docker service")
	}

	engine.Registry = &EnvironmentRegistry{
		Docker: dockr,
	}

	engine.Executor = &TaskExecutor{
		Docker: dockr,
	}

	return nil
}
