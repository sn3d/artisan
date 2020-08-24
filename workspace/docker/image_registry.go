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

// Build builds the docker image for faction. If image already exist, the function
// just returns you existing faction image.
func (cb *ImageRegistry) Build(fact *api.Faction, srcDir string) (*api.Image, error) {

	fmt.Printf(" - '%s': checking...", fact.Name)

	imageID := cb.GetImageID(fact)
	if imageID != "" {
		img := &api.Image{
			ID: imageID,
		}
		fmt.Printf("\033[2K\r - '%s': [OK]\n", fact.Name)
		return img, nil
	}

	var img *api.Image
	var err error

	if fact.Src != "" {
		fmt.Printf("\033[2K\r - '%s': building...\n", fact.Name)
		img, err = buildImage(cb.Docker, fact.Name, srcDir)
	} else {
		fmt.Printf("\033[2K\r - '%s': pulling...\n", fact.Name)
		img, err = pullImage(cb.Docker, fact.Image)
	}

	if err != nil {
		fmt.Printf("\033[0A\033[2K\r - '%s'': [ERROR]\n", fact.Name)
	} else {
		fmt.Printf("\033[0A\033[2K\r - '%s': [OK]\n", fact.Name)
	}

	return img, err
}

func (ir *ImageRegistry) GetImageID(fact *api.Faction) string {
	if fact.Image != "" {
		return getImageID(ir.Docker, fact.Image)
	} else {
		return getImageID(ir.Docker, factionToTag(fact.Name))
	}
}
