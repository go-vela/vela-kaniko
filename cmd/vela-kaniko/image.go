// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Image represents the plugin configuration for image information.
type Image struct {
	// variables passed to the image at build-time
	Args []string
	// path to the context for building the image
	Context string
	// path to the file for building the image
	Dockerfile string
	// build stage to target for image
	Target string
	// enable force adding metadata layers to build image
	ForceBuildMetadata bool
	// custom platform for image
	CustomPlatform string
}

// Validate verifies the Image is properly configured.
func (i *Image) Validate() error {
	logrus.Trace("validating image plugin configuration")

	// verify context is provided
	if len(i.Context) == 0 {
		return fmt.Errorf("no image context provided")
	}

	// verify dockerfile is provided
	if len(i.Dockerfile) == 0 {
		return fmt.Errorf("no image dockerfile provided")
	}

	return nil
}
