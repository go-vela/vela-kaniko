// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestDocker_Build_Validate(t *testing.T) {
	// setup types
	b := &Build{
		Event:        "push",
		Sha:          "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		Tag:          "v0.0.0",
		SnapshotMode: "redo",
	}

	err := b.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestDocker_Build_Validate_NoEvent(t *testing.T) {
	// setup types
	b := &Build{
		Sha: "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d",
		Tag: "v0.0.0",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Build_Validate_NoSha(t *testing.T) {
	// setup types
	b := &Build{
		Event: "push",
		Tag:   "v0.0.0",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestDocker_Build_Validate_InvalidSnapshotMode(t *testing.T) {
	// setup types
	b := &Build{
		Event:        "push",
		Tag:          "v0.0.0",
		SnapshotMode: "redo",
	}

	err := b.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
