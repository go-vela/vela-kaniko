// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

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

// helper function to validate fields being passed by user to plugin
func validate(e *env, p *plugin) error {

	// if the auto tag flag is set auto tag with commit or tag
	if p.AutoTag {
		switch e.BuildEvent {
		case "tag":
			p.Tags = append(p.Tags, e.BuildTag)
		default:
			p.Tags = append(p.Tags, e.BuildCommit)
		}
	}

	if len(p.Registry) == 0 {
		return fmt.Errorf("Plugin field registry is mandatory")
	}
	if len(p.Repo) == 0 {
		return fmt.Errorf("Plugin field repo is mandatory")
	}
	if len(p.Tags) == 0 {
		re := regexp.MustCompile(`^[A-Za-z0-9\-\.\_]*$`)

		// check each tag value for valid docker tag syntax
		// See docker docs for examples: https://docs.docker.com/engine/reference/commandline/tag/#Extended%20description
		for _, tag := range p.Tags {
			if !re.MatchString(tag) {
				return fmt.Errorf("Plugin tag %s has unaccepted rune. (Valid char set: abcdefghijklmnopqrstuvwxyz0123456789_-.ABCDEFGHIJKLMNOPQRSTUVWXY", tag)
			}
		}
		return fmt.Errorf("Plugin field tags or auto tag is mandatory")
	}
	if !strings.Contains("panic|fatal|error|warn|info|debug", p.LogLevel) {
		return fmt.Errorf("Plugin field debug accepted values are: panic|fatal|error|warn|info|debug")
	}

	if len(p.CacheRepo) != 0 {
		if !p.Cache {
			return fmt.Errorf("Plugin field cache is mandatory when using cache repo")
		}
	}

	if !p.DryRun {
		if len(p.Username) == 0 {
			return fmt.Errorf("Plugin field username or dry_run must be set")
		}

		if len(p.Password) == 0 {
			return fmt.Errorf("Plugin field password or dry_run must be set")
		}
	}

	return nil
}

// helper function to write a docker conf to kaniko user dir
func dockerConf(p plugin) error {

	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", p.Username, p.Password)))

	err := ioutil.WriteFile(fmt.Sprint("/kaniko/.docker/config.json"), []byte(fmt.Sprintf(conf, p.Registry, auth)), 0644)
	if err != nil {
		return err
	}

	return nil
}

// helper function to convert the plugin configuration into kaniko CLI flags
func buildCommand(p plugin) []string {

	flags := []string{}

	// get the working directory for context
	dir, err := os.Getwd()
	if err != nil {
		return nil
	}

	// Add required and default fields to kaniko command
	logrus.Debug("...add required kaniko flags")
	for _, tag := range p.Tags {
		flags = append(flags, fmt.Sprintf("--destination=%s:%s", p.Repo, tag))
	}

	flags = append(flags, fmt.Sprintf("--verbosity=%s", p.LogLevel))
	if len(p.Context) != 0 {
		flags = append(flags, fmt.Sprintf("--context=%s/%s", dir, p.Context))
	} else if dir != "/" {
		flags = append(flags, fmt.Sprintf("--context=%s", dir))
	}

	// handle adding optional fields to kaniko commands
	logrus.Debug("...add optional kaniko flags")
	if len(p.Dockerfile) != 0 {
		flags = append(flags, fmt.Sprintf("--dockerfile=%s", p.Dockerfile))
	}
	if p.DryRun {
		flags = append(flags, fmt.Sprint("--no-push"))
	}
	if len(p.BuildArgs) != 0 {
		for _, arg := range p.BuildArgs {
			flags = append(flags, fmt.Sprintf("--build-arg=%s", arg))
		}
	}
	if p.Cache {
		flags = append(flags, fmt.Sprint("--cache"))
		if len(p.CacheRepo) != 0 {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.CacheRepo))
		} else {
			flags = append(flags, fmt.Sprintf("--cache-repo=%s", p.Repo))
		}
	}

	logrus.Debugf("...flags added %+v", flags)

	return flags
}

// helper function to run the kaniko binary against provided plugin configuration
func run(flags []string) error {

	cmd := exec.Command("/kaniko/executor", flags...)
	var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Error: %s", err)
	}

	if errStdout != nil || errStderr != nil {
		return fmt.Errorf("Error: %s", err)
	}

	return nil
}
