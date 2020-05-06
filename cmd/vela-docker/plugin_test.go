// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"os/exec"
	"reflect"
	"testing"

	"github.com/spf13/afero"
)

func TestDocker_Plugin_Exec_BadWrite(t *testing.T) {
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
			Target:     "",
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
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDocker_Plugin_Exec_BadExec(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

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
			Target:     "",
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
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestDocker_Plugin_Command(t *testing.T) {
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
			Target:     "foo",
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

	want := exec.Command(
		kanikoBin,
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:v0.0.0",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--target=foo",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_NoCacheRepo(t *testing.T) {
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
			Target:     "",
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

	want := exec.Command(
		kanikoBin,
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		"--dockerfile=Dockerfile",
		"--no-push",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
	}
}

func TestDocker_Plugin_Command_NoDryRun(t *testing.T) {
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
			Target:     "",
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

	want := exec.Command(
		kanikoBin,
		"--build-arg=foo=bar",
		"--cache",
		"--cache-repo=index.docker.io/target/vela-docker",
		"--context=.",
		"--destination=index.docker.io/target/vela-docker:latest",
		"--destination=index.docker.io/target/vela-docker:7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		"--dockerfile=Dockerfile",
		"--verbosity=info",
	)

	// run test
	got := p.Command()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Command is %v, want %v", got, want)
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
			Target:     "",
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
			Target:     "",
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
			Target:     "",
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
			Target:     "",
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
