package docker

import (
	"fmt"
	"github.com/docker/docker/client"
	"github.com/unravela/delvin/api"
)

type ImageRegistry struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

// Build builds the docker image for class. If image already exist, the function
// just returns you existing class image.
func (cb *ImageRegistry) Build(cls *api.Class, srcDir string) (*api.Image, error) {

	fmt.Printf(" - '%s': checking...", cls.Name)

	imageID := cb.GetImageID(cls)
	if imageID != "" {
		img := &api.Image{
			ID: imageID,
		}
		fmt.Printf("\033[2K\r - '%s': [OK]\n", cls.Name)
		return img, nil
	}

	var img *api.Image
	var err error

	if cls.Src != "" {
		fmt.Printf("\033[2K\r - '%s': building...\n", cls.Name)
		img, err = buildImage(cb.Docker, cls.Name, srcDir)
	} else {
		fmt.Printf("\033[2K\r - '%s': pulling...\n", cls.Name)
		img, err = pullImage(cb.Docker, cls.Image)
	}

	if err != nil {
		fmt.Printf("\033[0A\033[2K\r - '%s'': [ERROR]\n", cls.Name)
	} else {
		fmt.Printf("\033[0A\033[2K\r - '%s': [OK]\n", cls.Name)
	}

	return img, err
}

func (ir *ImageRegistry) GetImageID(cls *api.Class) string {
	if cls.Image != "" {
		return getImageID(ir.Docker, cls.Image)
	} else {
		return getImageID(ir.Docker, classNameToTag(cls.Name))
	}
}
