//+build integration

package docker_test

import (
	"github.com/unravela/delvin/delvin"
	"github.com/unravela/delvin/delvin/docker"
	"testing"
)

// Requirements:
//  - This test require docker service with no image 'dlvin-jdk8' present
func TestClassBuilder_Build(t *testing.T) {
	// given a class with source defined
	ws, _ := delvin.Open("../../testdata/simplerepo")
	class := ws.Class("jdk8")
	if class == nil {
		t.Errorf("Cannot get the class")
	}

	// when we want to build class image
	cb, _ := docker.NewEnvClassBuilder()
	image, _ := cb.Build(class, ws)

	// then we should obtain some image ID
	if image.ID == "" {
		t.Errorf("Class image is empty string")
	}
}
