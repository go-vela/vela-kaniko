package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Build represents the plugin configuration for build information.
type Build struct {
	// event generated for build
	Event string
	// SHA-1 hash generated for commit
	Sha string
	// tag generated for build
	Tag string
}

// Validate verifies the Build is properly configured.
func (b *Build) Validate() error {
	logrus.Trace("validating build plugin configuration")

	if len(b.Event) == 0 {
		return fmt.Errorf("no build event provided")
	}

	if len(b.Sha) == 0 {
		return fmt.Errorf("no build sha provided")
	}

	return nil
}
