// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestDocker_Repo_Validate(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     true,
		CacheName: "index.docker.io/target/vela-kaniko",
		Name:      "index.docker.io/target/vela-kaniko",
		Tags:      []string{"latest"},
		AutoTag:   true,
		Label:     &Label{},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Repo_Validate_NoCache(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     false,
		CacheName: "index.docker.io/target/vela-kaniko",
		Name:      "index.docker.io/target/vela-kaniko",
		Tags:      []string{"latest"},
		AutoTag:   true,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_NoName(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     true,
		CacheName: "index.docker.io/target/vela-kaniko",
		Name:      "",
		Tags:      []string{"latest"},
		AutoTag:   true,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_InvalidTags(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     true,
		CacheName: "",
		Name:      "index.docker.io/target/vela-kaniko",
		Tags:      []string{"!@#$%^&*()"},
		AutoTag:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_NoTags(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     true,
		CacheName: "",
		Name:      "index.docker.io/target/vela-kaniko",
		Tags:      []string{},
		AutoTag:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
