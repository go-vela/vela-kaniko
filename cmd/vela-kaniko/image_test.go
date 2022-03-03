// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestDocker_Image_Validate(t *testing.T) {
	// setup types
	i := &Image{
		Args:               []string{},
		Context:            ".",
		Dockerfile:         "Dockerfile",
		Target:             "",
		ForceBuildMetadata: false,
	}

	err := i.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Image_Validate_NoContext(t *testing.T) {
	// setup types
	i := &Image{
		Args:       []string{},
		Dockerfile: "Dockerfile",
		Target:     "",
	}

	err := i.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Image_Validate_NoDockerfile(t *testing.T) {
	// setup types
	i := &Image{
		Args:    []string{},
		Context: ".",
		Target:  "",
	}

	err := i.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
