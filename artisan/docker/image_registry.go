package docker

import (
	"fmt"

	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
)

// ImageRegistry implementation acessing to docker images
type ImageRegistry struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

// Build builds the docker image for environment. If image already exist, the function
// just returns you existing env. image.
func (ir *ImageRegistry) Build(env *api.Environment, srcDir string) (*api.Image, error) {

	fmt.Printf(" - '%s': checking...", env.Name)

	imageID := ir.GetImageID(env)
	if imageID != "" {
		img := &api.Image{
			ID: imageID,
		}
		fmt.Printf("\033[2K\r - '%s': [OK]\n", env.Name)
		return img, nil
	}

	var img *api.Image
	var err error

	if env.Src != "" {
		fmt.Printf("\033[2K\r - '%s': building...\n", env.Name)
		img, err = buildImage(ir.Docker, env.Name, srcDir)
	} else {
		fmt.Printf("\033[2K\r - '%s': pulling...\n", env.Name)
		img, err = pullImage(ir.Docker, env.Image)
	}

	if err != nil {
		fmt.Printf("\033[0A\033[2K\r - '%s'': [ERROR]\n", env.Name)
	} else {
		fmt.Printf("\033[0A\033[2K\r - '%s': [OK]\n", env.Name)
	}

	return img, err
}

// GetImageID returns you docker image ID for given environment
func (ir *ImageRegistry) GetImageID(env *api.Environment) string {
	if env.Image != "" {
		return getImageID(ir.Docker, env.Image)
	}
	return getImageID(ir.Docker, envToTag(env.Name))
}
