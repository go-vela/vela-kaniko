// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"regexp"
	"strings"

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
		// type of compression - 'gzip' (default) or 'zstd'
		Compression string
		// level of compression - 1 to 19 (inclusive), default is -1
		CompressionLevel int
		// used for translating the pre-defined image labels
		Label *Label
		// labels of the image for the repository
		Labels []string
		// name of the repository for the image
		Name string
		// tags of the image for the repository
		Tags []string
		// a filter for topics
		TopicsFilter string
	}

	// Label represents the open image specification fields.
	// Note: these are populated by Vela build variables.
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
		// topics attached to repository
		Topics []string
		// direct url of the repository
		URL string
	}
)

// AddLabels container labels to the image.
func (r *Repo) AddLabels() []string {
	// store all the topics
	topics := r.Label.Topics
	// labels we will return
	labels := []string{}

	// append the standard set of labels
	labels = append(
		labels,
		fmt.Sprintf("org.opencontainers.image.created=%s", r.Label.Created),
		fmt.Sprintf("org.opencontainers.image.url=%s", r.Label.URL),
		fmt.Sprintf("org.opencontainers.image.revision=%s", r.Label.Commit),
		fmt.Sprintf("io.vela.build.author=%s", r.Label.AuthorEmail),
		fmt.Sprintf("io.vela.build.number=%d", r.Label.Number),
		fmt.Sprintf("io.vela.build.repo=%s", r.Label.FullName),
		fmt.Sprintf("io.vela.build.commit=%s", r.Label.Commit),
		fmt.Sprintf("io.vela.build.url=%s", r.Label.URL),
	)

	// if a filter is defined, use it to only
	// include those that match the filter
	if len(r.TopicsFilter) > 0 {
		// clear out current topics
		topics = []string{}
		// we already confirmed validity of regex expression in
		// .Validate, so we skip the error check here
		re, _ := regexp.Compile(r.TopicsFilter)
		for _, topic := range r.Label.Topics {
			if re.MatchString(topic) {
				topics = append(topics, topic)
			}
		}
	}

	// only append topics if we have any, since the possibility of
	// it being empty is much higher than any of the other fields
	// that are derived from standard Vela build variables
	if len(topics) > 0 {
		labels = append(labels, fmt.Sprintf("io.vela.build.topics=%s", strings.Join(topics, ",")))
	}

	return labels
}

// ConfigureAutoTagBuildTags adds the build tag to repo tags.
func (r *Repo) ConfigureAutoTagBuildTags(b *Build) {
	// check what build event was provided
	switch b.Event {
	case "tag":
		// add build tag to list of repo tags
		r.Tags = append(r.Tags, b.Tag)
	default:
		// add build sha to list of repo tags
		r.Tags = append(r.Tags, b.Sha)
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
		// check each tag value for valid docker tag syntax
		for _, tag := range r.Tags {
			if !tagRegexp.MatchString(tag) {
				return fmt.Errorf(errTagValidation, tag)
			}
		}
	}

	// check validity of regex expression for topics
	if len(r.TopicsFilter) > 0 {
		_, err := regexp.Compile(r.TopicsFilter)
		if err != nil {
			return fmt.Errorf("topics filter regex not valid")
		}
	}

	// make sure a valid compression type was provided, if any
	if len(r.Compression) > 0 {
		if r.Compression != "gzip" && r.Compression != "zstd" {
			return fmt.Errorf("compression has to be one of 'gzip' or 'zstd'")
		}
	}

	// make sure compression level is between 1 and 19
	if r.CompressionLevel > 0 && r.CompressionLevel > 19 {
		return fmt.Errorf("compression-level can't exceed 19")
	}

	return nil
}
