package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
)

// EnvironmentRegistry implementation for docker resolve the
// env. definitions to docker images.
type EnvironmentRegistry struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

// Build builds the docker image for environment. If image already exist, the function
// just returns you existing env. image.
func (ir *EnvironmentRegistry) Build(envDef *api.EnvironmentDef, srcDir string) (api.EnvironmentID, error) {

	fmt.Printf(" - '%s': checking...", envDef.Name)

	imageID := ir.getDockerImageID(envDef)
	if imageID != "" {
		img := api.EnvironmentID(imageID)
		fmt.Printf("\033[2K\r - '%s': [OK]\n", envDef.Name)
		return img, nil
	}

	var img api.EnvironmentID
	var err error

	if envDef.Src != "" {
		fmt.Printf("\033[2K\r - '%s': building...\n", envDef.Name)
		img, err = buildImage(ir.Docker, envDef.Name, srcDir)
	} else {
		fmt.Printf("\033[2K\r - '%s': pulling...\n", envDef.Name)
		img, err = pullImage(ir.Docker, envDef.Image)
	}

	if err != nil {
		fmt.Printf("\033[0A\033[2K\r - '%s'': [ERROR]\n", envDef.Name)
	} else {
		fmt.Printf("\033[0A\033[2K\r - '%s': [OK]\n", envDef.Name)
	}

	return img, err
}

// getDockerImageID returns you docker image ID for given environment definition.
func (ir *EnvironmentRegistry) getDockerImageID(envDef *api.EnvironmentDef) string {
	if envDef.Image != "" {
		return getImageID(ir.Docker, envDef.Image)
	}
	return getImageID(ir.Docker, envToTag(envDef.Name))
}
