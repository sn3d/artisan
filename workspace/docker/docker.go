package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/unravela/artisan/api"
)

// ImagePrefix is used in all docker images generated from source. If
// faction is named as "node", the image will be "artsn-node:latest"
const ImagePrefix = "artsn-"

func pullImage(docker *client.Client, image string) (*api.Image, error) {

	if image == "" {
		return nil, errors.New("no image is present")
	}

	// We need canonical name of the image. That means not only something:latest but
	// also registry need to be part of the string. For now,
	// we are using docker.io/library as default
	img := "docker.io/library/" + image

	ctx := context.Background()
	res, err := docker.ImagePull(ctx, img, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot pull faction image %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Printf("\033[2K\r    > %s", txt)
	}

	imageID := getImageID(docker, image)
	out := &api.Image{
		ID: imageID,
	}
	return out, nil
}

// this function builds docker image for given source dir. The image is build
// with given name
func buildImage(docker *client.Client, name string, srcDir string) (*api.Image, error) {
	ctx := context.Background()
	// create tar
	if _, err := os.Stat(srcDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("invalid path to source dir")
	}

	bctx := createBuildContext(buildContextOptions{
		root:     srcDir,
		includes: []string{},
	})

	// build image
	tags := []string{factionToTag(name)}
	res, err := docker.ImageBuild(ctx, bctx, types.ImageBuildOptions{
		Dockerfile: "./Dockerfile",
		NoCache:    true,
		Tags:       tags,
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot build faction %w", err)
	}

	// print the image build result to output
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		txt := scanner.Text()
		fmt.Printf("\033[2K\r    > %s", txt)
	}

	imageID := getImageID(docker, factionToTag(name))
	out := &api.Image{
		ID: imageID,
	}

	return out, nil
}

// function returns you docker image ID for faction.
// This is needed e.g. when you want to perform task and task is
// running in forge.
//
// If function returns you empty string, that means there is no
// docker image present in system for this forge and forge need to
// be build.
func getImageID(docker *client.Client, image string) string {
	ctx := context.Background()
	images, _ := docker.ImageList(ctx, types.ImageListOptions{})
	for _, img := range images {
		tags := img.RepoTags

		idx := len(tags) - 1
		if idx < 0 {
			idx = 0
		}

		for _, tag := range tags {
			if tag == image {
				return img.ID
			}
		}
	}

	return ""
}

// This function transform faction name e.g. '@jdk8' to
// docker tag 'dlvin-jdk8:latest'
func factionToTag(name string) string {
	if strings.HasPrefix(name, "@") {
		name = name[1:]
	}
	return ImagePrefix + name + ":latest"
}
