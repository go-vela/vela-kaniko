package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Registry represents the plugin configuration for registry information.
//
// https://docs.docker.com/registry/
type Registry struct {
	// name of the registry to publish the image to
	Name string
	// user name for communication with the registry
	Username string
	// password for communication with the registry
	Password string
	// enable building the image without publishing
	DryRun bool
}

// Validate verifies the Registry is properly configured.
func (r *Registry) Validate() error {
	logrus.Trace("validating registry plugin configuration")

	if len(r.Name) == 0 {
		return fmt.Errorf("no registry name provided")
	}

	if !r.DryRun {
		if len(r.Username) == 0 {
			return fmt.Errorf("no registry username provided")
		}

		if len(r.Password) == 0 {
			return fmt.Errorf("no registry password provided")
		}
	}

	return nil
}
