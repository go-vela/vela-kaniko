// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/spf13/afero"
)

func TestDocker_Registry_Validate(t *testing.T) {
	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Registry_Validate_NoName(t *testing.T) {
	// setup types
	r := &Registry{
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Registry_Validate_NoUsername(t *testing.T) {
	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Registry_Validate_NoPassword(t *testing.T) {
	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		DryRun:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Registry_Write(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoName(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Username: "octocat",
		Password: "superSecretPassword",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoUsername(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}

func TestDocker_Registry_Write_NoPassword(t *testing.T) {
	// setup filesystem
	appFS = afero.NewMemMapFs()

	// setup types
	r := &Registry{
		Name:     "index.docker.io",
		Username: "octocat",
		DryRun:   false,
	}

	err := r.Write()
	if err != nil {
		t.Errorf("Write returned err: %v", err)
	}
}
