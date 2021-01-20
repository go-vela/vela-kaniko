// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"regexp"

	"github.com/sirupsen/logrus"
)

type (
	// Repo represents the plugin configuration for repo information.
	Repo struct {
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
		// used for translating the pre-defined image labels
		Label *Label
		// labels of the image for the repository
		Labels []string
	}

	// Label represents the open image specification fields.
	Label struct {
		// author from the source commit
		AuthorEmail string
		// commit sha from the source commit
		Commit string
		// timestamp when the image was built
		Created string
		// full name of the repository
		FullName string
		// build number from vela
		Number int
		// direct url of the repository
		URL string
	}
)

// AddLabels adds open container spec labels to plugin
//
// https://github.com/opencontainers/image-spec/blob/v1.0.1/annotations.md
func (r *Repo) AddLabels() []string {
	return []string{
		fmt.Sprintf("org.opencontainers.image.created=%s", r.Label.Created),
		fmt.Sprintf("org.opencontainers.image.url=%s", r.Label.URL),
		fmt.Sprintf("org.opencontainers.image.revision=%s", r.Label.Commit),
		fmt.Sprintf("io.vela.build.author=%s", r.Label.AuthorEmail),
		fmt.Sprintf("io.vela.build.number=%d", r.Label.Number),
		fmt.Sprintf("io.vela.build.repo=%s", r.Label.FullName),
		fmt.Sprintf("io.vela.build.commit=%s", r.Label.Commit),
		fmt.Sprintf("io.vela.build.url=%s", r.Label.URL),
	}
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

		// nolint
		// check each tag value for valid docker tag syntax
		// See docker docs for examples: https://docs.docker.com/engine/reference/commandline/tag/#Extended%20description
		for _, tag := range r.Tags {
			if !re.MatchString(tag) {
				return fmt.Errorf("invalid tag provided: %s (Valid char set: abcdefghijklmnopqrstuvwxyz0123456789_-.ABCDEFGHIJKLMNOPQRSTUVWXY", tag)
			}
		}
	}

	// add pre-defined labels
	r.Labels = append(r.Labels, r.AddLabels()...)

	return nil
}
