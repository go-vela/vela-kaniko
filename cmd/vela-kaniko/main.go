// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-vela/vela-kaniko/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
)

// nolint: funlen // ignore function length due to comments and flags
func main() {
	// capture application version information
	v := version.New()

	// serialize the version information as pretty JSON
	bytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	// output the version information to stdout
	fmt.Fprintf(os.Stdout, "%s\n", string(bytes))

	// create new CLI application
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-kaniko"
	app.HelpName = "vela-kaniko"
	app.Usage = "Vela Kaniko plugin for building and publishing images"
	app.Copyright = "Copyright (c) 2021 Target Brands, Inc. All rights reserved."
	app.Authors = []*cli.Author{
		{
			Name:  "Vela Admins",
			Email: "vela@target.com",
		},
	}

	// Plugin Metadata

	app.Action = run
	app.Compiled = time.Now()
	app.Version = v.Semantic()

	// Plugin Flags

	// nolint
	app.Flags = []cli.Flag{

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_LOG_LEVEL", "KANIKO_LOG_LEVEL"},
			FilePath: "/vela/parameters/kaniko/log_level,/vela/secrets/kaniko/log_level",
			Name:     "log.level",
			Usage:    "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Value:    "info",
		},

		// Build Flags

		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_EVENT", "KANIKO_EVENT", "VELA_BUILD_EVENT"},
			FilePath: "/vela/parameters/kaniko/event,/vela/secrets/kaniko/event",
			Name:     "build.event",
			Usage:    "event triggered for build",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SHA", "KANIKO_SHA", "VELA_BUILD_COMMIT"},
			FilePath: "/vela/parameters/kaniko/sha,/vela/secrets/kaniko/sha",
			Name:     "build.sha",
			Usage:    "commit SHA-1 hash for build",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_SNAPSHOT_MODE", "KANIKO_SNAPSHOT_MODE"},
			FilePath: "/vela/parameters/kaniko/snapshot_mode,/vela/secrets/kaniko/snapshot_mode",
			Name:     "build.snapshot_mode",
			Usage:    "control how to snapshot the filesystem. - options (full|redo|time)",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TAG", "KANIKO_TAG", "VELA_BUILD_TAG"},
			FilePath: "/vela/parameters/kaniko/tag,/vela/secrets/kaniko/tag",
			Name:     "build.tag",
			Usage:    "full tag reference for build (only populated for tag events)",
		},

		// Image Flags

		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_BUILD_ARGS", "KANIKO_BUILD_ARGS"},
			FilePath: "/vela/parameters/kaniko/build_args,/vela/secrets/kaniko/build_args",
			Name:     "image.build_args",
			Usage:    "variables passed to the image at build-time",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_CONTEXT", "KANIKO_CONTEXT"},
			FilePath: "/vela/parameters/kaniko/context,/vela/secrets/kaniko/context",
			Name:     "image.context",
			Usage:    "path on local filesystem for building image from",
			Value:    ".",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_DOCKERFILE", "KANIKO_DOCKERFILE"},
			FilePath: "/vela/parameters/kaniko/dockerfile,/vela/secrets/kaniko/dockerfile",
			Name:     "image.dockerfile",
			Usage:    "path to text file with build instructions",
			Value:    "Dockerfile",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TARGET", "KANIKO_TARGET"},
			FilePath: "/vela/parameters/kaniko/target,/vela/secrets/kaniko/target",
			Name:     "image.target",
			Usage:    "build stage to target for image",
		},

		// Registry Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_DRY_RUN", "KANIKO_DRY_RUN"},
			FilePath: "/vela/parameters/kaniko/dry_run,/vela/secrets/kaniko/dry_run",
			Name:     "registry.dry_run",
			Usage:    "enables building images without publishing to the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REGISTRY", "KANIKO_REGISTRY"},
			FilePath: "/vela/parameters/kaniko/registry,/vela/secrets/kaniko/registry",
			Name:     "registry.name",
			Usage:    "Docker registry name to communicate with",
			Value:    "index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_MIRROR", "KANIKO_MIRROR"},
			FilePath: "/vela/parameters/kaniko/mirror,/vela/secrets/kaniko/mirror",
			Name:     "registry.mirror",
			Usage:    "name of the mirror registry to use instead of index.docker.io",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_USERNAME", "KANIKO_USERNAME", "DOCKER_USERNAME"},
			FilePath: "/vela/parameters/kaniko/username,/vela/secrets/kaniko/username",
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "KANIKO_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: "/vela/parameters/kaniko/password,/vela/secrets/kaniko/password",
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_PUSH_RETRY", "KANIKO_PUSH_RETRY"},
			FilePath: "/vela/parameters/kaniko/push_retry,/vela/secrets/kaniko/push_retry",
			Name:     "registry.push_retry",
			Usage:    "number of retries for pushing an image to a remote destination",
		},

		// Repo Flags

		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_AUTO_TAG", "KANIKO_AUTO_TAG"},
			FilePath: "/vela/parameters/kaniko/auto_tag,/vela/secrets/kaniko/auto_tag",
			Name:     "repo.auto_tag",
			Usage:    "enables automatically providing tags for the image",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_CACHE", "KANIKO_CACHE"},
			FilePath: "/vela/parameters/kaniko/cache,/vela/secrets/kaniko/cache",
			Name:     "repo.cache",
			Usage:    "enables caching of each layer for the image",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_CACHE_REPO", "KANIKO_CACHE_REPO"},
			FilePath: "/vela/parameters/kaniko/cache_repo,/vela/secrets/kaniko/cache_repo",
			Name:     "repo.cache_name",
			Usage:    "enables caching of each layer for a specific repo for the image",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REPO", "KANIKO_REPO"},
			FilePath: "/vela/parameters/kaniko/repo,/vela/secrets/kaniko/repo",
			Name:     "repo.name",
			Usage:    "repository name for the image",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_TAGS", "KANIKO_TAGS"},
			FilePath: "/vela/parameters/kaniko/tags,/vela/secrets/kaniko/tags",
			Name:     "repo.tags",
			Usage:    "repository tags of the image",
			Value:    cli.NewStringSlice("latest"),
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_LABELS", "KANIKO_LABELS"},
			FilePath: "/vela/parameters/kaniko/labels,/vela/secrets/kaniko/labels",
			Name:     "repo.labels",
			Usage:    "repository labels of the image",
		},

		// extract vars for open image specification labeling
		&cli.StringFlag{
			EnvVars: []string{"VELA_BUILD_AUTHOR_EMAIL"},
			Name:    "label.author_email",
			Usage:   "author from the source commit",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_BUILD_COMMIT"},
			Name:    "label.commit",
			Usage:   "commit sha from the source commit",
		},
		&cli.IntFlag{
			EnvVars: []string{"VELA_BUILD_NUMBER"},
			Name:    "label.number",
			Usage:   "build number",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_REPO_FULL_NAME"},
			Name:    "label.full_name",
			Usage:   "full name of the repository",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_REPO_LINK"},
			Name:    "label.url",
			Usage:   "direct url of the repository",
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		logrus.Fatal(err)
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
		"code":     "https://github.com/go-vela/vela-kaniko",
		"docs":     "https://go-vela.github.io/docs/plugins/registry/pipeline/kaniko",
		"registry": "https://hub.docker.com/r/target/vela-kaniko",
	}).Info("Vela Kaniko Plugin")

	// create the plugin
	p := &Plugin{
		// build configuration
		Build: &Build{
			Event:        c.String("build.event"),
			Sha:          c.String("build.sha"),
			SnapshotMode: c.String("build.snapshot_mode"),
			Tag:          c.String("build.tag"),
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
			DryRun:    c.Bool("registry.dry_run"),
			Name:      c.String("registry.name"),
			Mirror:    c.String("registry.mirror"),
			Username:  c.String("registry.username"),
			Password:  c.String("registry.password"),
			PushRetry: c.Int("registry.push_retry"),
		},
		// repo configuration
		Repo: &Repo{
			AutoTag:   c.Bool("repo.auto_tag"),
			Cache:     c.Bool("repo.cache"),
			CacheName: c.String("repo.cache_name"),
			Name:      c.String("repo.name"),
			Tags:      c.StringSlice("repo.tags"),
			Label: &Label{
				AuthorEmail: c.String("label.author_email"),
				Commit:      c.String("label.commit"),
				Created:     time.Now().Format(time.RFC3339),
				FullName:    c.String("label.full_name"),
				Number:      c.Int("label.number"),
				URL:         c.String("label.url"),
			},
			Labels: c.StringSlice("repo.labels"),
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
