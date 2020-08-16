//+build integration

package docker

import (
	"github.com/unravela/delvin/delvin"
	"testing"
)

// This test require docker service with no image
func TestClassBuilder_Build(t *testing.T) {
	// given a class with source defined
	ws, err := delvin.Open("../../testdata/simplerepo")
	if err != nil {
		t.Errorf("Cannot open workspace %w", err)
	}

	class := ws.Class("jdk8")
	if class == nil {
		t.Errorf("Cannot get the class")
	}

	// when we want to build class image
	cb, err := NewEnvClassBuilder()
	if err != nil {
		t.Errorf("Error creating class builder: %w", err)
	}

	image, err := cb.Build(class, ws)
	if err != nil {
		t.Errorf("")
	}

	// then we should obtain the image's ID
	if image.ID == "" {
		t.Errorf("Class image is empty string")
	}
}
