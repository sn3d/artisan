package docker

import (
	"github.com/docker/docker/client"
	"github.com/unravela/delvin/delvin"
)

type ClassBuilder struct {
	// ref to docker client that is used for all docker operations
	Docker *client.Client
}

// NewEnvClassBuilder initialize class builder with docker client based on
// docker env. variables like DOCKER_HOST etc..
func NewEnvClassBuilder() (*ClassBuilder, error) {
	var cbuilder *ClassBuilder
	var err error

	if dockr, err := client.NewEnvClient(); err == nil {
		cbuilder = &ClassBuilder{
			Docker: dockr,
		}
	}

	return cbuilder, err
}

// Build builds the docker image for class. If image already exist, the function
// just returns you existing class image.
func (cb *ClassBuilder) Build(cls *delvin.Class, ws *delvin.Workspace) (*delvin.ClassImage, error) {
	imageID := cb.GetImageID(cls)
	if imageID != "" {
		img := &delvin.ClassImage{
			ID: imageID,
		}
		return img, nil
	}

	if cls.Src != "" {
		srcDir := ws.AbsPath(cls.Src)
		return buildImage(cb.Docker, cls.Name, srcDir)
	} else {
		return pullImage(cb.Docker, cls.Image)
	}
}

func (cb *ClassBuilder) GetImageID(cls *delvin.Class) string {
	if cls.Image != "" {
		return getImageID(cb.Docker, cls.Image)
	} else {
		return getImageID(cb.Docker, classNameToTag(cls.Name))
	}
}