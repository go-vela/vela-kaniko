// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"reflect"
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
			Args:       []string{"foo=bar"},
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

func TestDocker_Plugin_Flags(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "tag",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   true,
		},
		Repo: &Repo{
			Cache:     true,
			CacheName: "index.docker.io/target/vela-docker",
			Name:      "index.docker.io/target/vela-docker",
			Tags:      []string{"latest"},
			AutoTag:   true,
		},
	}

	want := []string{
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:v0.0.0",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--verbosity=info",
	}

	// run test
	got := p.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flags is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Flags_NoCacheRepo(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
			Context:    ".",
			Dockerfile: "Dockerfile",
		},
		Registry: &Registry{
			Name:     "index.docker.io",
			Username: "octocat",
			Password: "superSecretPassword",
			DryRun:   true,
		},
		Repo: &Repo{
			Cache:   true,
			Name:    "index.docker.io/target/vela-docker",
			Tags:    []string{"latest"},
			AutoTag: true,
		},
	}

	want := []string{
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--verbosity=info",
	}

	// run test
	got := p.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flags is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Flags_NoDryRun(t *testing.T) {
	// setup types
	p := &Plugin{
		Build: &Build{
			Event: "push",
			Sha:   "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
			Tag:   "v0.0.0",
		},
		Image: &Image{
			Args:       []string{"foo=bar"},
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

	want := []string{
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		"--dockerfile=Dockerfile",
		"--verbosity=info",
	}

	// run test
	got := p.Flags()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flags is %v, want %v", got, want)
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
			Args:       []string{"foo=bar"},
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
			Args:       []string{"foo=bar"},
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
			Args:       []string{"foo=bar"},
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
			Args:       []string{"foo=bar"},
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
