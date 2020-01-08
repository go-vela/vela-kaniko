// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-docker"
	app.HelpName = "vela-docker"
	app.Usage = "Vela Docker plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2019 Target Brands, Inc. All rights reserved."
	app.Authors = []cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Compiled = time.Now()
	app.Action = run

	// Plugin Flags

	app.Flags = []cli.Flag{

		cli.StringFlag{
			EnvVar: "PARAMETER_LOG_LEVEL,VELA_LOG_LEVEL,DOCKER_LOG_LEVEL",
			Name:   "log.level",
			Usage:  "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:  "info",
		},

		// Build Flags

		cli.StringFlag{
			EnvVar: "PARAMETER_EVENT,BUILD_EVENT",
			Name:   "build.event",
			Usage:  "event triggered for build",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_SHA,BUILD_COMMIT",
			Name:   "build.sha",
			Usage:  "commit SHA-1 hash for build",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_TAG,BUILD_TAG",
			Name:   "build.tag",
			Usage:  "full tag reference for build (only populated for tag events)",
		},

		// Image Flags

		cli.StringSliceFlag{
			EnvVar: "PARAMETER_BUILD_ARGS,IMAGE_BUILD_ARGS",
			Name:   "image.build_args",
			Usage:  "variables passed to the image at build-time",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_CONTEXT,IMAGE_CONTEXT",
			Name:   "image.context",
			Usage:  "path on local filesystem for building image from",
			Value:  ".",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_DOCKERFILE,IMAGE_DOCKERFILE",
			Name:   "image.dockerfile",
			Usage:  "path to text file with build instructions",
			Value:  "Dockerfile",
		},

		// Registry Flags

		cli.BoolFlag{
			EnvVar: "PARAMETER_DRY_RUN,REGISTRY_DRY_RUN",
			Name:   "registry.dry_run",
			Usage:  "enables building images without publishing to the registry",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_REGISTRY,REGISTRY_NAME",
			Name:   "registry.name",
			Usage:  "Docker registry name to communicate with",
			Value:  "index.docker.io",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_USERNAME,REGISTRY_USERNAME,DOCKER_USERNAME",
			Name:   "registry.username",
			Usage:  "user name for communication with the registry",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_PASSWORD,REGISTRY_PASSWORD,DOCKER_PASSWORD",
			Name:   "registry.password",
			Usage:  "password for communication with the registry",
		},

		// Repo Flags

		cli.BoolFlag{
			EnvVar: "PARAMETER_AUTO_TAG,REPO_AUTO_TAG",
			Name:   "repo.auto_tag",
			Usage:  "enables automatically providing tags for the image",
		},
		cli.BoolFlag{
			EnvVar: "PARAMETER_CACHE,REPO_CACHE",
			Name:   "repo.cache",
			Usage:  "enables caching of each layer for the image",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_CACHE_REPO,REPO_CACHE_NAME",
			Name:   "repo.cache_name",
			Usage:  "enables caching of each layer for a specific repo for the image",
		},
		cli.StringFlag{
			EnvVar: "PARAMETER_REPO,REPO_NAME",
			Name:   "repo.name",
			Usage:  "repository name for the image",
		},
		cli.StringSliceFlag{
			EnvVar:   "PARAMETER_TAGS,REPO_TAGS",
			FilePath: ".tags",
			Name:     "repo.tags",
			Usage:    "repository tags of the image",
			Value:    &cli.StringSlice{"latest"},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(c *cli.Context) error {
	// set the log level for the plugin
	switch c.String("log.level") {
	case "t", "trace", "Trace", "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
	case "d", "debug", "Debug", "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "w", "warn", "Warn", "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "e", "error", "Error", "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "f", "fatal", "Fatal", "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	case "p", "panic", "Panic", "PANIC":
		logrus.SetLevel(logrus.PanicLevel)
	case "i", "info", "Info", "INFO":
		fallthrough
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	// create the plugin
	p := &Plugin{
		// build configuration
		Build: &Build{
			Event: c.String("build.event"),
			Sha:   c.String("build.sha"),
			Tag:   c.String("build.tag"),
		},
		// image configuration
		Image: &Image{
			Args:       c.StringSlice("image.build_args"),
			Context:    c.String("image.context"),
			Dockerfile: c.String("image.dockerfile"),
		},
		// registry configuration
		Registry: &Registry{
			DryRun:   c.Bool("registry.dry_run"),
			Name:     c.String("registry.name"),
			Username: c.String("registry.username"),
			Password: c.String("registry.password"),
		},
		// repo configuration
		Repo: &Repo{
			AutoTag:   c.Bool("repo.auto_tag"),
			Cache:     c.Bool("repo.cache"),
			CacheName: c.String("repo.cache_name"),
			Name:      c.String("repo.name"),
			Tags:      c.StringSlice("repo.tags"),
		},
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec()
}
