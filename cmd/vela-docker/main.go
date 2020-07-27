// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-docker"
	app.HelpName = "vela-docker"
	app.Usage = "Vela Docker plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2020 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
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

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "VELA_LOG_LEVEL", "DOCKER_LOG_LEVEL"},
			FilePath: string("/vela/secrets/parameter/log_level,/vela/secrets/docker/log_level"),
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Build Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_EVENT", "BUILD_EVENT"},
			FilePath: string("/vela/secrets/parameter/build/event,/vela/secrets/docker/build/event"),
			Name:     "build.event",
			Usage:    "event triggered for build",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SHA", "BUILD_COMMIT"},
			FilePath: string("/vela/secrets/parameter/build/sha,/vela/secrets/docker/build/sha"),
			Name:     "build.sha",
			Usage:    "commit SHA-1 hash for build",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TAG", "BUILD_TAG"},
			FilePath: string("/vela/secrets/parameter/build/tag,/vela/secrets/docker/build/tag"),
			Name:     "build.tag",
			Usage:    "full tag reference for build (only populated for tag events)",
		},

		// Image Flags

		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_BUILD_ARGS", "IMAGE_BUILD_ARGS"},
			FilePath: string("/vela/secrets/parameter/image/build_args,/vela/secrets/docker/image/build_args"),
			Name:     "image.build_args",
			Usage:    "variables passed to the image at build-time",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_CONTEXT", "IMAGE_CONTEXT"},
			FilePath: string("/vela/secrets/parameter/image/context,/vela/secrets/docker/image/context"),
			Name:     "image.context",
			Usage:    "path on local filesystem for building image from",
			Value:    ".",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_DOCKERFILE", "IMAGE_DOCKERFILE"},
			FilePath: string("/vela/secrets/parameter/image/dockerfile,/vela/secrets/docker/image/dockerfile"),
			Name:     "image.dockerfile",
			Usage:    "path to text file with build instructions",
			Value:    "Dockerfile",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET", "IMAGE_TARGET"},
			FilePath: string("/vela/secrets/parameter/image/target,/vela/secrets/docker/image/target"),
			Name:     "image.target",
			Usage:    "build stage to target for image",
		},

		// Registry Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "REGISTRY_DRY_RUN"},
			FilePath: string("/vela/secrets/parameter/dry_run,/vela/secrets/docker/dry_run"),
			Name:     "registry.dry_run",
			Usage:    "enables building images without publishing to the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "REGISTRY_NAME"},
			FilePath: string("/vela/secrets/parameter/registry/name,/vela/secrets/docker/registry/name"),
			Name:     "registry.name",
			Usage:    "Docker registry name to communicate with",
			Value:    "index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "REGISTRY_USERNAME", "DOCKER_USERNAME"},
			FilePath: string("/vela/secrets/parameter/registry/username,/vela/secrets/docker/registry/username,/vela/secrets/docker/username"),
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "REGISTRY_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: string("/vela/secrets/parameter/registry/password,/vela/secrets/docker/registry/password,/vela/secrets/docker/password"),
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},

		// Repo Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_AUTO_TAG", "REPO_AUTO_TAG"},
			FilePath: string("/vela/secrets/parameter/repo/auto_tag,/vela/secrets/docker/repo/auto_tag"),
			Name:     "repo.auto_tag",
			Usage:    "enables automatically providing tags for the image",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_CACHE", "REPO_CACHE"},
			FilePath: string("/vela/secrets/parameter/repo/cache,/vela/secrets/docker/repo/cache"),
			Name:     "repo.cache",
			Usage:    "enables caching of each layer for the image",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_CACHE_REPO", "REPO_CACHE_NAME"},
			FilePath: string("/vela/secrets/parameter/repo/cache_name,/vela/secrets/docker/repo/cache_name"),
			Name:     "repo.cache_name",
			Usage:    "enables caching of each layer for a specific repo for the image",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REPO", "REPO_NAME"},
			FilePath: string("/vela/secrets/parameter/repo/name,/vela/secrets/docker/repo/name"),
			Name:     "repo.name",
			Usage:    "repository name for the image",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_TAGS", "REPO_TAGS"},
			FilePath: string(".tags,/vela/secrets/parameter/repo/tags,/vela/secrets/docker/repo/tags"),
			Name:     "repo.tags",
			Usage:    "repository tags of the image",
			Value:    cli.NewStringSlice("latest"),
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_LABELS", "REPO_LABELS"},
			FilePath: string("/vela/secrets/parameter/repo/labels,/vela/secrets/docker/repo/labels"),
			Name:     "repo.labels",
			Usage:    "repository labels of the image",
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

	logrus.WithFields(logrus.Fields{
		"code":     "https://github.com/go-vela/vela-docker",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/docker",
		"registry": "https://hub.docker.com/r/target/vela-docker",
	}).Info("Vela Docker Plugin")

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
			Target:     c.String("image.target"),
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
			Labels:    c.StringSlice("repo.labels"),
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
