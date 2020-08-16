package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/unravela/delvin/delvin"
	"os"
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

func (cb *ClassBuilder) Build(cls *delvin.Class, ws *delvin.Workspace) (*delvin.ClassImage, error) {
	// we rebuild only when docker image isn't present
	imageID := cb.getImageID(cls)
	if imageID != "" {
		img := &delvin.ClassImage{
			ID: imageID,
		}
		return img, nil
	}

	if cls.Src != "" {
		return cb.buildImage(cls, ws)
	} else {
		return cb.pullImage(cls)
	}
}

// function returns you docker image ID for class.
// This is needed e.g. when you want to perform task and task is
// running in forge.
//
// If function returns you empty string, that means there is no
// docker image present in system for this forge and forge need to
// be build.
func (cb *ClassBuilder) getImageID(cls *delvin.Class) string {

	ctx := context.Background()
	images, _ := cb.Docker.ImageList(ctx, types.ImageListOptions{})
	for _, img := range images {
		tags := img.RepoTags

		// we compare last tag - because there are also parent's tags
		// this would be great to redesign and search in all tags
		idx := len(tags) - 1
		if idx < 0 {
			idx = 0
		}

		if cls.Image == "" {
			for _, tag := range tags {
				if tag == "ignt-"+cls.Name+":latest" {
					return img.ID
				}
			}
		} else {
			if tags[idx] == cls.Image {
				return img.ID
			}
		}
	}

	return ""
}

// this function executes docker image building for forge's source.
// The building is executed only when Src is set. The forge's name is
// also used as repo-tag for docker image.
func (cb *ClassBuilder) buildImage(cls *delvin.Class, ws *delvin.Workspace) (*delvin.ClassImage, error) {
	ctx := context.Background()
	srcDir := ws.AbsPath(cls.Src)
	// create tar
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		panic("Invalid path to forge")
	}

	bctx := createBuildContext(buildContextOptions{
		root:     srcDir,
		includes: []string{},
	})

	// build image
	tags := []string{"ignt-" + cls.Name}
	res, err := cb.Docker.ImageBuild(ctx, bctx, types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		NoCache:    true,
		Tags:       tags,
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot build class %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		println(scanner.Text())
	}

	imageId := cb.getImageID(cls)
	out := &delvin.ClassImage{
		ID: imageId,
	}

	return out, nil
}

func (cb *ClassBuilder) pullImage(cls *delvin.Class) (*delvin.ClassImage, error) {

	if cls.Image == "" {
		return nil, errors.New("No image is present!")
	}

	// pull the image
	ctx := context.Background()

	// We need canonical name of the image. That means not only something:latest but
	// also registry need to be part of the string. For now,
	// we are using docker.io/library as default
	img := "docker.io/library/" + cls.Image

	res, err := cb.Docker.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot pull class image %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	out := &delvin.ClassImage{
		ID: "",
	}
	return out, nil
}
