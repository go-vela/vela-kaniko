// SPDX-License-Identifier: Apache-2.0

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

func TestDocker_Repo_Validate_BadTopicsFilter(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:        true,
		CacheName:    "",
		Name:         "index.docker.io/target/vela-kaniko",
		Tags:         []string{},
		AutoTag:      true,
		TopicsFilter: "?()",
		Label: &Label{
			Topics: []string{},
		},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_Validate_ValidTopicsFilter(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:        true,
		CacheName:    "",
		Name:         "index.docker.io/target/vela-kaniko",
		Tags:         []string{},
		AutoTag:      true,
		TopicsFilter: "^id",
		Label: &Label{
			Topics: []string{"id123"},
		},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate should not have returned err")
	}
}

func TestDocker_Repo_Validate_SingleLabel(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:        true,
		CacheName:    "",
		Name:         "index.docker.io/target/vela-kaniko",
		Tags:         []string{},
		AutoTag:      true,
		TopicsFilter: "^id",
		Label: &Label{
			Topics: []string{"id123"},
		},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate should not have returned err")
	}
}

func TestDocker_Repo_Compression(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:       true,
		CacheName:   "index.docker.io/target/vela-kaniko",
		Compression: "zstd",
		Name:        "index.docker.io/target/vela-kaniko",
		Tags:        []string{"latest"},
		AutoTag:     true,
		Label:       &Label{},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Repo_BadCompression(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:       true,
		CacheName:   "index.docker.io/target/vela-kaniko",
		Compression: "foo",
		Name:        "index.docker.io/target/vela-kaniko",
		Tags:        []string{"latest"},
		AutoTag:     true,
		Label:       &Label{},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_CompressionLevel(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:            true,
		CacheName:        "index.docker.io/target/vela-kaniko",
		Compression:      "zstd",
		CompressionLevel: 2,
		Name:             "index.docker.io/target/vela-kaniko",
		Tags:             []string{"latest"},
		AutoTag:          true,
		Label:            &Label{},
	}

	err := r.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Repo_CompressionLevelTooHigh(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:            true,
		CacheName:        "index.docker.io/target/vela-kaniko",
		Compression:      "zstd",
		CompressionLevel: 50,
		Name:             "index.docker.io/target/vela-kaniko",
		Tags:             []string{"latest"},
		AutoTag:          true,
		Label:            &Label{},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Repo_CompressionLevelTooLow(t *testing.T) {
	// setup types
	r := &Repo{
		Cache:            true,
		CacheName:        "index.docker.io/target/vela-kaniko",
		Compression:      "zstd",
		CompressionLevel: -1,
		Name:             "index.docker.io/target/vela-kaniko",
		Tags:             []string{"latest"},
		AutoTag:          true,
		Label:            &Label{},
	}

	err := r.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
