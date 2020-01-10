// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"os/exec"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"
)

var appFS = afero.NewOsFs()

// Plugin represents the configuration loaded for the plugin.
type Plugin struct {
	// build arguments loaded for the plugin
	Build *Build
	// image arguments loaded for the plugin
	Image *Image
	// registry arguments loaded for the plugin
	Registry *Registry
	// repo arguments loaded for the plugin
	Repo *Repo
}

// Command formats and outputs the command necessary for
// Kaniko to build and publish a Docker image.
func (p *Plugin) Command() *exec.Cmd {
	var flags []string

	for _, arg := range p.Image.Args {
		flags = append(flags, fmt.Sprintf("--build-arg=%s", arg))
	}

	if p.Repo.Cache {
		flags = append(flags, fmt.Sprint("--cache"))

		if len(p.Repo.CacheName) > 0 {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.Repo.CacheName))
		} else {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.Repo.Name))
		}
	}

	flags = append(flags, fmt.Sprintf("--context=%s", p.Image.Context))

	if p.Repo.AutoTag {
		switch p.Build.Event {
		case "tag":
			p.Repo.Tags = append(p.Repo.Tags, p.Build.Tag)
		default:
			p.Repo.Tags = append(p.Repo.Tags, p.Build.Sha)
		}
	}

	for _, tag := range p.Repo.Tags {
		flags = append(flags, fmt.Sprintf("--destination=%s:%s", p.Repo.Name, tag))
	}

	flags = append(flags, fmt.Sprintf("--dockerfile=%s", p.Image.Dockerfile))

	if p.Registry.DryRun {
		flags = append(flags, fmt.Sprint("--no-push"))
	}

	flags = append(flags, fmt.Sprintf("--verbosity=%s", logrus.GetLevel()))

	return exec.Command(kanikoBin, flags...)
}

// Exec formats and runs the commands for building and publishing a Docker image.
func (p *Plugin) Exec() error {
	logrus.Debug("running plugin with provided configuration")

	err := p.Registry.Write()
	if err != nil {
		return err
	}

	err = execCmd(p.Command())
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Plugin is properly configured.
func (p *Plugin) Validate() error {
	logrus.Debug("validating plugin configuration")

	// validate build configuration
	err := p.Build.Validate()
	if err != nil {
		return err
	}

	// validate image configuration
	err = p.Image.Validate()
	if err != nil {
		return err
	}

	// validate registry configuration
	err = p.Registry.Validate()
	if err != nil {
		return err
	}

	// validate repo configuration
	err = p.Repo.Validate()
	if err != nil {
		return err
	}

	return nil
}