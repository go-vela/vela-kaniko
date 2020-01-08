// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

const (
	credentials = `%s:%s`

	registryFile = `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`
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

// Write creates a Docker config.json file for building and publishing the image.
func (r *Registry) Write() error {
	a := &afero.Afero{
		Fs: appFS,
	}

	if len(r.Name) == 0 || len(r.Username) == 0 || len(r.Password) == 0 {
		return nil
	}

	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, r.Username, r.Password)),
	)

	out := fmt.Sprintf(
		registryFile,
		r.Name,
		basicAuth,
	)

	path := "/kaniko/.docker/config.json"

	return a.WriteFile(path, []byte(out), 0644)
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
