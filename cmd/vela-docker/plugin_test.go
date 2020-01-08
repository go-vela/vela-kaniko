// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestDocker_Plugin_Exec(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
	}
}

func TestDocker_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Plugin_Validate_NoBuild(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{},
		Image: &Image{
			Args:       []string{},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoImage(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoRegistry(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Plugin_Validate_NoRepo(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   false,
		},
		Repo: &Repo{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
