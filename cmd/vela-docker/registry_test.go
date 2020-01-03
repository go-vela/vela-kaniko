package main

import "testing"

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
