// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
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
	// insecure registries to push/pull from
	InsecureRegistries []string
	// name of the mirror registry to use instead of index.docker.io
	Mirror string
	// name of the registry to publish the image to
	Name string
	// user name for communication with the registry
	Username string
	// password for communication with the registry
	Password string
	// enable building the image without publishing
	PushRetry int
	// enable pulling from any insecure registry
	DryRun bool
	// number of retries for pushing an image to a remote destination
	InsecurePull bool
	// enable pushing to any insecure registry
	InsecurePush bool
}

// Write creates a Docker config.json file for building and publishing the image.
func (r *Registry) Write() error {
	logrus.Trace("writing registry configuration file")

	// use custom filesystem which enables us to test
	a := &afero.Afero{
		Fs: appFS,
	}

	// check if name, username and password are provided
	if len(r.Name) == 0 || len(r.Username) == 0 || len(r.Password) == 0 {
		return nil
	}

	// create basic authentication string for config.json file
	basicAuth := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf(credentials, r.Username, r.Password)),
	)

	// create output string for config.json file
	out := fmt.Sprintf(
		registryFile,
		r.Name,
		basicAuth,
	)

	// create full path for config.json file
	path := "/kaniko/.docker/config.json"

	// nolint: gomnd // ignore magic number
	return a.WriteFile(path, []byte(out), 0644)
}

// Validate verifies the Registry is properly configured.
func (r *Registry) Validate() error {
	logrus.Trace("validating registry plugin configuration")

	// verify registry is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no registry name provided")
	}

	// check if dry run is disabled
	if !r.DryRun {
		// check if username is provided
		if len(r.Username) == 0 {
			return fmt.Errorf("no registry username provided")
		}

		// check if password is provided
		if len(r.Password) == 0 {
			return fmt.Errorf("no registry password provided")
		}
	}

	return nil
}
