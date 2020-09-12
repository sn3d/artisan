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

// Build builds the docker image for faction. If image already exist, the function
// just returns you existing faction image.
func (ir *ImageRegistry) Build(fact *api.Faction, srcDir string) (*api.Image, error) {

	fmt.Printf(" - '%s': checking...", fact.Name)

	imageID := ir.GetImageID(fact)
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
		img, err = buildImage(ir.Docker, fact.Name, srcDir)
	} else {
		fmt.Printf("\033[2K\r - '%s': pulling...\n", fact.Name)
		img, err = pullImage(ir.Docker, fact.Image)
	}

	if err != nil {
		fmt.Printf("\033[0A\033[2K\r - '%s'': [ERROR]\n", fact.Name)
	} else {
		fmt.Printf("\033[0A\033[2K\r - '%s': [OK]\n", fact.Name)
	}

	return img, err
}

// GetImageID returns you docker image ID for given faction
func (ir *ImageRegistry) GetImageID(fact *api.Faction) string {
	if fact.Image != "" {
		return getImageID(ir.Docker, fact.Image)
	}
	return getImageID(ir.Docker, factionToTag(fact.Name))
}
