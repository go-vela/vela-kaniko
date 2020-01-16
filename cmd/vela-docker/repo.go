// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

// Repo represents the plugin configuration for repo information.
type Repo struct {
	// enable tagging of image via commit or tag
	AutoTag bool
	// enable caching of image layers
	Cache bool
	// enable caching of image layers for a specific repo
	CacheName string
	// name of the repository for the image
	Name string
	// tags of the image for the repository
	Tags []string
}

// Validate verifies the Repo is properly configured.
func (r *Repo) Validate() error {
	logrus.Trace("validating repo plugin configuration")

	// check if cache name is provided
	if len(r.CacheName) > 0 {
		// verify caching is enabled
		if !r.Cache {
			return fmt.Errorf("cache not set for cache repo: %s", r.CacheName)
		}
	}

	// verify repo is provided
	if len(r.Name) == 0 {
		return fmt.Errorf("no repo name provided")
	}

	// check if auto tagging is disabled
	if !r.AutoTag {
		// verify tags are provided
		if len(r.Tags) == 0 {
			return fmt.Errorf("no repo tags provided")
		}
	}

	// check if tags are provided
	if len(r.Tags) > 0 {
		// create regular expression for verifying tags
		re := regexp.MustCompile(`^[A-Za-z0-9\-\.\_]*$`)

		// check each tag value for valid docker tag syntax
		// See docker docs for examples: https://docs.docker.com/engine/reference/commandline/tag/#Extended%20description
		for _, tag := range r.Tags {
			if !re.MatchString(tag) {
				return fmt.Errorf("invalid tag provided: %s (Valid char set: abcdefghijklmnopqrstuvwxyz0123456789_-.ABCDEFGHIJKLMNOPQRSTUVWXY", tag)
			}
		}
	}

	return nil
}
