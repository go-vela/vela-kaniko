// SPDX-License-Identifier: Apache-2.0

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
		// type of compression - 'gzip' (default if not defined) or 'zstd'
		Compression string
		// level of compression - 1 to 9 (inclusive)
		CompressionLevel int
		// prevent tar compression for cached layers
		CompressedCaching bool
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
		// direct url of the build
		BuildURL string
		// host that the image is built on
		Host string
		// custom set of labels that the user can provide
		CustomSet []string
	}
)

// AddLabels container labels to the image.
func (r *Repo) AddLabels() []string {
	// store all the topics
	topics := r.Label.Topics

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

	labelMap := map[string]string{
		"org.opencontainers.image.created":  r.Label.Created,
		"org.opencontainers.image.url":      r.Label.URL,
		"org.opencontainers.image.revision": r.Label.Commit,
		"io.vela.build.author":              r.Label.AuthorEmail,
		"io.vela.build.number":              fmt.Sprintf("%d", r.Label.Number),
		"io.vela.build.repo":                r.Label.FullName,
		"io.vela.build.commit":              r.Label.Commit,
		"io.vela.build.url":                 r.Label.URL,
		"io.vela.build.link":                r.Label.BuildURL,
		"io.vela.build.host":                r.Label.Host,
	}

	if len(topics) > 0 {
		labelMap["io.vela.build.topics"] = strings.Join(topics, ",")
	}

	// labels we will return
	labels := []string{}

	// append the standard set of labels
	for k, v := range labelMap {
		labels = append(labels, fmt.Sprintf("%s=%s", k, v))
	}

	// append the custom set of labels
	for _, label := range r.Label.CustomSet {
		parsed := strings.Split(label, "=")

		// do not let user overwrite predefined labels
		if _, exists := labelMap[parsed[0]]; exists {
			logrus.Fatalf("custom label %s already exists in predefined labels", parsed[0])
		}

		labels = append(labels, label)
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
			return fmt.Errorf("compression must be one of 'gzip' or 'zstd'")
		}
	}

	// make sure compression level is between 1 and 9 inclusive
	if r.CompressionLevel != 0 {
		if r.CompressionLevel < 1 || r.CompressionLevel > 9 {
			return fmt.Errorf("compression-level must be between 1 - 9 inclusive")
		}
	}

	if len(r.Label.CustomSet) > 0 {
		for _, label := range r.Label.CustomSet {
			split := strings.Split(label, "=")

			if len(split) != 2 {
				return fmt.Errorf("custom label %s is not in the format key=value", label)
			}
		}
	}

	return nil
}
