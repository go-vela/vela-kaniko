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
}

// Validate verifies the Image is properly configured.
func (i *Image) Validate() error {
	logrus.Trace("validating image plugin configuration")

	if len(i.Context) == 0 {
		return fmt.Errorf("no image context provided")
	}

	if len(i.Dockerfile) == 0 {
		return fmt.Errorf("no image dockerfile provided")
	}

	return nil
}
