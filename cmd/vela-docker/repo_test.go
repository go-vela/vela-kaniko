package main

import "testing"

func TestDocker_Repo_Validate(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:     true,
		CacheName: "index.docker.io/target/vela-docker",
		Name:      "index.docker.io/target/vela-docker",
		Tags:      []string{"latest"},
		AutoTag:   true,
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
		CacheName: "index.docker.io/target/vela-docker",
		Name:      "index.docker.io/target/vela-docker",
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
		CacheName: "index.docker.io/target/vela-docker",
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
		Name:      "index.docker.io/target/vela-docker",
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
		Name:      "index.docker.io/target/vela-docker",
		Tags:      []string{},
		AutoTag:   false,
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
