// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// plugin struct represents fields user can present to plugin
type plugin struct {
	Registry   string   // registry you plan to upload docker image
	Repo       string   // repository name for the image
	Tags       []string // repository tag for the image
	Username   string   // authenticates with this username
	Password   string   // authenticates with this password
	Dockerfile string   // dockerfile to be used, defaults to Dockerfile
	DryRun     bool     // boolean if the docker image should be pushed at the end
	Context    string   // the context path to use, defaults to root of the git repo
	BuildArgs  []string // custom arguments passed to docker build
	LogLevel   string   // enable verbose logs on plugin
	Cache      bool     // enable docker image layer caching
	CacheRepo  string   // enable docker image layer caching for a specific repo, note cache=true is reqired with this flag
	AutoTag    bool     // enable plugin to auto tag the docker image via commit or tag
}

// env struct represents the environment variables the CI gives you for free
type env struct {
	BuildEvent  string
	BuildCommit string
	BuildTag    string
}

const conf = `{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`

func main() {
	app := cli.NewApp()
	app.Name = "docker"
	app.Usage = "Docker plugin for building docker images"
	app.Action = setup
	app.Flags = []cli.Flag{

		// Required flags
		cli.StringFlag{
			Name:   "registry",
			EnvVar: "PARAMETER_REGISTRY,DOCKER_REGISTRY",
		},
		cli.StringFlag{
			Name:   "repo",
			EnvVar: "PARAMETER_REPO,DOCKER_REPO",
		},
		cli.StringSliceFlag{
			Name:   "tags",
			EnvVar: "PARAMETER_TAGS,DOCKER_TAGS",
		},

		// Authentication Flags
		cli.StringFlag{
			Name:   "username",
			EnvVar: "PARAMETER_USERNAME,DOCKER_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			EnvVar: "PARAMETER_PASSWORD,DOCKER_PASSWORD",
		},

		// Optional flags
		cli.StringFlag{
			Name:   "dockerfile",
			EnvVar: "PARAMETER_DOCKERFILE,DOCKER_DOCKERFILE",
		},
		cli.BoolFlag{
			Name:   "dry-run",
			EnvVar: "PARAMETER_DRY_RUN,DOCKER_DRY_RUN",
		},
		cli.StringFlag{
			Name:   "context",
			EnvVar: "PARAMETER_CONTEXT,DOCKER_CONTEXT",
		},
		cli.StringSliceFlag{
			Name:   "build-args",
			EnvVar: "PARAMETER_BUILD_ARGS,DOCKER_BUILD_ARGS",
		},
		cli.StringFlag{
			Name:   "log-level",
			Usage:  "valid values: panic|fatal|error|warn|info|debug",
			EnvVar: "PARAMETER_LOG_LEVEL,DOCKER_LOG_LEVEL",
			Value:  "info",
		},
		cli.BoolFlag{
			Name:   "cache",
			EnvVar: "PARAMETER_CACHE,DOCKER_CACHE",
		},
		cli.StringFlag{
			Name:   "cache-repo",
			EnvVar: "PARAMETER_CACHE_REPO,DOCKER_CACHE_REPO",
		},
		cli.BoolFlag{
			Name:   "auto-tag",
			EnvVar: "PARAMETER_AUTO_TAG,DOCKER_AUTO_TAG",
		},

		// These fields are passed into the environment via the default environment variables
		cli.StringFlag{
			Name:   "build-event",
			EnvVar: "BUILD_EVENT",
		},
		cli.StringFlag{
			Name:   "build-commit",
			EnvVar: "BUILD_COMMIT",
		},
		cli.StringFlag{
			Name:   "build-tag",
			EnvVar: "BUILD_TAG",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// helper function used as entrypoint of plugin execution
func setup(c *cli.Context) error {

	// set application log information
	if c.String("debug") != "info" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// read plugin from context into plugin struct
	p := plugin{
		Registry:   c.String("registry"),
		Repo:       c.String("repo"),
		Tags:       c.StringSlice("tags"),
		Username:   c.String("username"),
		Password:   c.String("password"),
		Dockerfile: c.String("dockerfile"),
		DryRun:     c.Bool("dry-run"),
		Context:    c.String("context"),
		BuildArgs:  c.StringSlice("build-args"),
		LogLevel:   c.String("log-level"),
		Cache:      c.Bool("cache"),
		CacheRepo:  c.String("cache-repo"),
		AutoTag:    c.Bool("auto-tag"),
	}

	// read environment from context into env struct
	env := env{
		BuildEvent:  c.String("build-event"),
		BuildCommit: c.String("build-commit"),
		BuildTag:    c.String("build-tag"),
	}

	// check plugin configuration has proper fields
	logrus.Info("Validating plugin configuration....")
	err := validate(&env, &p)
	if err != nil {
		return err
	}
	logrus.Info("Configuration is valid...")

	// Write docker config to kaniko directory with registry credentials
	logrus.Info("Creating docker config with credentials...")
	err = dockerConf(p)
	if err != nil {
		return err
	}
	logrus.Info("Credentials created...")

	// convert the plugin configuration into kaniko CLI flags
	logrus.Info("Build Kaniko flags from provided plugin configuration...")
	flags := buildCommand(p)
	logrus.Info("Kaniko flags built...")

	logrus.Info("Execute kaniko docker plugin...")
	err = run(flags)
	if err != nil {
		return err
	}
	logrus.Info("Plugin finished...")

	return nil
}
