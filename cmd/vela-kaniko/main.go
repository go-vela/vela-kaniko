// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-vela/vela-kaniko/version"
)

//nolint:funlen // ignore function length due to comments and flags
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
	app := &cli.Command{
		Name:      "vela-kaniko",
		Usage:     "Vela Kaniko plugin for building and publishing images",
		Copyright: "Copyright 2019 Target Brands, Inc. All rights reserved.",
		Authors: []any{
			&mail.Address{
				Name:    "Vela Admins",
				Address: "vela@target.com",
			},
		},
		// Plugin Metadata
		Version: v.Semantic(),
		Action:  run,
	}

	// Plugin Flags
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "log.level",
			Value: "info",
			Usage: "set log level - options: (trace|debug|info|warn|error|fatal|panic)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_LEVEL"),
				cli.EnvVar("KANIKO_LOG_LEVEL"),
				cli.File("/vela/parameters/kaniko/log_level"),
				cli.File("/vela/secrets/kaniko/log_level"),
			),
		},

		// Build Flags
		&cli.StringFlag{
			Name:  "build.event",
			Usage: "event triggered for build",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_EVENT"),
				cli.EnvVar("KANIKO_EVENT"),
				cli.EnvVar("VELA_BUILD_EVENT"),
				cli.File("/vela/parameters/kaniko/event"),
				cli.File("/vela/secrets/kaniko/event"),
			),
		},
		&cli.StringFlag{
			Name:  "build.sha",
			Usage: "commit SHA-1 hash for build",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SHA"),
				cli.EnvVar("KANIKO_SHA"),
				cli.EnvVar("VELA_BUILD_COMMIT"),
				cli.File("/vela/parameters/kaniko/sha"),
				cli.File("/vela/secrets/kaniko/sha"),
			),
		},
		&cli.StringFlag{
			Name:  "build.snapshot_mode",
			Usage: "control how to snapshot the filesystem - options (full|redo|time)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SNAPSHOT_MODE"),
				cli.EnvVar("KANIKO_SNAPSHOT_MODE"),
				cli.File("/vela/parameters/kaniko/snapshot_mode"),
				cli.File("/vela/secrets/kaniko/snapshot_mode"),
			),
		},
		&cli.StringFlag{
			Name:  "build.tag",
			Usage: "full tag reference for build (only populated for tag events)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TAG"),
				cli.EnvVar("KANIKO_TAG"),
				cli.EnvVar("VELA_BUILD_TAG"),
				cli.File("/vela/parameters/kaniko/tag"),
				cli.File("/vela/secrets/kaniko/tag"),
			),
		},
		&cli.BoolFlag{
			Name:  "build.use_new_run",
			Usage: "use the experimental run implementation for detecting changes without requiring file system snapshots - in some cases, this may improve build performance by 75%",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_USE_NEW_RUN"),
				cli.EnvVar("KANIKO_USE_NEW_RUN"),
				cli.EnvVar("VELA_BUILD_USE_NEW_RUN"),
				cli.File("/vela/parameters/kaniko/use_new_run"),
				cli.File("/vela/secrets/kaniko/use_new_run"),
			),
		},
		&cli.StringFlag{
			Name:  "build.tar_path",
			Usage: "if set, the image will be saved as a tarball at that path",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TAR_PATH"),
				cli.EnvVar("KANIKO_TAR_PATH"),
				cli.EnvVar("VELA_BUILD_TAR_PATH"),
				cli.File("/vela/parameters/kaniko/tar_path"),
				cli.File("/vela/secrets/kaniko/tar_path"),
			),
		},
		&cli.BoolFlag{
			Name:  "build.single_snapshot",
			Usage: "takes a single snapshot of the filesystem at the end of the build, so only one layer will be appended to the base image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_SINGLE_SNAPSHOT"),
				cli.EnvVar("KANIKO_SINGLE_SNAPSHOT"),
				cli.EnvVar("VELA_BUILD_SINGLE_SNAPSHOT"),
				cli.File("/vela/parameters/kaniko/single_snapshot"),
				cli.File("/vela/secrets/kaniko/single_snapshot"),
			),
		},
		&cli.BoolFlag{
			Name:  "build.ignore_var_run",
			Value: true,
			Usage: "by default kaniko ignores /var/run when taking image snapshot - include this parameter to preserve /var/run/* in destination image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_IGNORE_VAR_RUN"),
				cli.EnvVar("KANIKO_IGNORE_VAR_RUN"),
				cli.EnvVar("VELA_IGNORE_VAR_RUN"),
				cli.File("/vela/parameters/kaniko/ignore_var_run"),
				cli.File("/vela/secrets/kaniko/ignore_var_run"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "build.ignore_path",
			Usage: "ignore paths when taking image snapshot",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_IGNORE_PATH"),
				cli.EnvVar("KANIKO_IGNORE_PATH"),
				cli.File("/vela/parameters/kaniko/ignore_path"),
				cli.File("/vela/secrets/kaniko/ignore_path"),
			),
		},
		&cli.BoolFlag{
			Name:  "build.log_timestamps",
			Usage: "enable timestamps in logs",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LOG_TIMESTAMPS"),
				cli.EnvVar("KANIKO_LOG_TIMESTAMPS"),
				cli.File("/vela/parameters/kaniko/log_timestamps"),
				cli.File("/vela/secrets/kaniko/log_timestamps"),
			),
		},

		// Image Flags
		&cli.StringFlag{
			Name:  "image.build_args",
			Usage: "variables passed to the image at build-time",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_BUILD_ARGS"),
				cli.EnvVar("KANIKO_BUILD_ARGS"),
				cli.File("/vela/parameters/kaniko/build_args"),
				cli.File("/vela/secrets/kaniko/build_args"),
			),
		},
		&cli.StringFlag{
			Name:  "image.context",
			Value: ".",
			Usage: "path on local filesystem for building image from",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_CONTEXT"),
				cli.EnvVar("KANIKO_CONTEXT"),
				cli.File("/vela/parameters/kaniko/context"),
				cli.File("/vela/secrets/kaniko/context"),
			),
		},
		&cli.StringFlag{
			Name:  "image.dockerfile",
			Value: "Dockerfile",
			Usage: "path to text file with build instructions",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_DOCKERFILE"),
				cli.EnvVar("KANIKO_DOCKERFILE"),
				cli.File("/vela/parameters/kaniko/dockerfile"),
				cli.File("/vela/secrets/kaniko/dockerfile"),
			),
		},
		&cli.StringFlag{
			Name:  "image.target",
			Usage: "build stage to target for image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TARGET"),
				cli.EnvVar("KANIKO_TARGET"),
				cli.File("/vela/parameters/kaniko/target"),
				cli.File("/vela/secrets/kaniko/target"),
			),
		},
		&cli.StringFlag{
			Name:  "image.force_build_metadata",
			Usage: "enables force adding metadata layers to build image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_FORCE_BUILD_METADATA"),
				cli.EnvVar("KANIKO_FORCE_BUILD_METADATA"),
				cli.File("/vela/parameters/kaniko/force_build_metadata"),
				cli.File("/vela/secrets/kaniko/force_build_metadata"),
			),
		},
		&cli.StringFlag{
			Name:  "image.custom_platform",
			Usage: "custom platform for the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_CUSTOM_PLATFORM"),
				cli.EnvVar("KANIKO_CUSTOM_PLATFORM"),
				cli.File("/vela/parameters/kaniko/custom_platform"),
				cli.File("/vela/secrets/kaniko/custom_platform"),
			),
		},

		// Registry Flags
		&cli.BoolFlag{
			Name:  "registry.dry_run",
			Usage: "enables building images without publishing to the registry",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_DRY_RUN"),
				cli.EnvVar("KANIKO_DRY_RUN"),
				cli.File("/vela/parameters/kaniko/dry_run"),
				cli.File("/vela/secrets/kaniko/dry_run"),
			),
		},
		&cli.StringFlag{
			Name:  "registry.name",
			Value: "index.docker.io",
			Usage: "docker registry name to communicate with",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REGISTRY"),
				cli.EnvVar("KANIKO_REGISTRY"),
				cli.File("/vela/parameters/kaniko/registry"),
				cli.File("/vela/secrets/kaniko/registry"),
			),
		},
		&cli.StringFlag{
			Name:  "registry.mirror",
			Usage: "name of the mirror registry to use instead of index.docker.io",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_MIRROR"),
				cli.EnvVar("KANIKO_MIRROR"),
				cli.File("/vela/parameters/kaniko/mirror"),
				cli.File("/vela/secrets/kaniko/mirror"),
			),
		},
		&cli.StringFlag{
			Name:  "registry.username",
			Usage: "user name for communication with the registry",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_USERNAME"),
				cli.EnvVar("KANIKO_USERNAME"),
				cli.EnvVar("DOCKER_USERNAME"),
				cli.File("/vela/parameters/kaniko/username"),
				cli.File("/vela/secrets/kaniko/username"),
				cli.File("/vela/secrets/managed-auth/username"),
			),
		},
		&cli.StringFlag{
			Name:  "registry.password",
			Usage: "password for communication with the registry",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_PASSWORD"),
				cli.EnvVar("KANIKO_PASSWORD"),
				cli.EnvVar("DOCKER_PASSWORD"),
				cli.File("/vela/parameters/kaniko/password"),
				cli.File("/vela/secrets/kaniko/password"),
				cli.File("/vela/secrets/managed-auth/password"),
			),
		},
		&cli.IntFlag{
			Name:  "registry.push_retry",
			Usage: "number of retries for pushing an image to a remote destination",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_PUSH_RETRY"),
				cli.EnvVar("KANIKO_PUSH_RETRY"),
				cli.File("/vela/parameters/kaniko/push_retry"),
				cli.File("/vela/secrets/kaniko/push_retry"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "registry.insecure_registries",
			Usage: "insecure registries to push & pull from",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("VELA_REGISTRY_INSECURE_REGISTRIES"),
				cli.EnvVar("KANIKO_INSECURE_REGISTRIES"),
				cli.File("/vela/parameters/kaniko/insecure_registries"),
				cli.File("/vela/secrets/kaniko/insecure_registries"),
			),
		},
		&cli.BoolFlag{
			Name:  "registry.insecure_pull",
			Usage: "enable pulling from insecure registries",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_INSECURE_PULL"),
				cli.EnvVar("KANIKO_INSECURE_PULL"),
				cli.File("/vela/parameters/kaniko/insecure_pull"),
				cli.File("/vela/secrets/kaniko/insecure_pull"),
			),
		},
		&cli.BoolFlag{
			Name:  "registry.insecure_push",
			Usage: "enable pushing to insecure registries",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_INSECURE_PUSH"),
				cli.EnvVar("KANIKO_INSECURE_PUSH"),
				cli.File("/vela/parameters/kaniko/insecure_push"),
				cli.File("/vela/secrets/kaniko/insecure_push"),
			),
		},

		// Repo Flags
		&cli.BoolFlag{
			Name:  "repo.auto_tag",
			Usage: "enables automatically providing tags for the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_AUTO_TAG"),
				cli.EnvVar("KANIKO_AUTO_TAG"),
				cli.File("/vela/parameters/kaniko/auto_tag"),
				cli.File("/vela/secrets/kaniko/auto_tag"),
			),
		},
		&cli.BoolFlag{
			Name:  "repo.cache",
			Usage: "enables caching of each layer for the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_CACHE"),
				cli.EnvVar("KANIKO_CACHE"),
				cli.File("/vela/parameters/kaniko/cache"),
				cli.File("/vela/secrets/kaniko/cache"),
			),
		},
		&cli.StringFlag{
			Name:  "repo.cache_name",
			Usage: "enables caching of each layer for a specific repo for the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_CACHE_REPO"),
				cli.EnvVar("KANIKO_CACHE_REPO"),
				cli.File("/vela/parameters/kaniko/cache_repo"),
				cli.File("/vela/secrets/kaniko/cache_repo"),
			),
		},
		&cli.StringFlag{
			Name:  "repo.compression",
			Usage: "set the compression type - gzip (default) or zstd",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_COMPRESSION"),
				cli.EnvVar("KANIKO_COMPRESSION"),
				cli.File("/vela/parameters/kaniko/compression"),
				cli.File("/vela/secrets/kaniko/compression"),
			),
		},
		&cli.IntFlag{
			Name:  "repo.compression_level",
			Usage: "set the compression level (1-9, inclusive)",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_COMPRESSION_LEVEL"),
				cli.EnvVar("KANIKO_COMPRESSION_LEVEL"),
				cli.File("/vela/parameters/kaniko/compression_level"),
				cli.File("/vela/secrets/kaniko/compression_level"),
			),
		},
		&cli.BoolFlag{
			Name:  "repo.compressed_caching",
			Value: true,
			Usage: "when set to false, will prevent tar compression for cached layers",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_COMPRESSED_CACHING"),
				cli.EnvVar("KANIKO_COMPRESSED_CACHING"),
				cli.File("/vela/parameters/kaniko/compressed_caching"),
				cli.File("/vela/secrets/kaniko/compressed_caching"),
			),
		},
		&cli.StringFlag{
			Name:  "repo.name",
			Usage: "repository name for the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REPO"),
				cli.EnvVar("KANIKO_REPO"),
				cli.File("/vela/parameters/kaniko/repo"),
				cli.File("/vela/secrets/kaniko/repo"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "repo.tags",
			Value: []string{"latest"},
			Usage: "repository tags of the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_TAGS"),
				cli.EnvVar("KANIKO_TAGS"),
				cli.File("/vela/parameters/kaniko/tags"),
				cli.File("/vela/secrets/kaniko/tags"),
			),
		},
		&cli.StringSliceFlag{
			Name:  "repo.labels",
			Usage: "repository labels of the image",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_LABELS"),
				cli.EnvVar("KANIKO_LABELS"),
				cli.File("/vela/parameters/kaniko/labels"),
				cli.File("/vela/secrets/kaniko/labels"),
			),
		},
		&cli.StringFlag{
			Name:  "repo.topics_filter",
			Usage: "filter to restrict which repository topics to include in label",
			Sources: cli.NewValueSourceChain(
				cli.EnvVar("PARAMETER_REPO_TOPICS_FILTER"),
				cli.EnvVar("KANIKO_REPO_TOPICS_FILTER"),
				cli.File("/vela/parameters/kaniko/repo_topics_filter"),
				cli.File("/vela/secrets/kaniko/repo_topics_filter"),
			),
		},

		// extract vars for open image specification labeling
		&cli.StringFlag{
			Name:    "label.author_email",
			Usage:   "author from the source commit",
			Sources: cli.EnvVars("VELA_BUILD_AUTHOR_EMAIL"),
		},
		&cli.StringFlag{
			Name:    "label.commit",
			Usage:   "commit sha from the source commit",
			Sources: cli.EnvVars("VELA_BUILD_COMMIT"),
		},
		&cli.IntFlag{
			Name:    "label.number",
			Usage:   "build number",
			Sources: cli.EnvVars("VELA_BUILD_NUMBER"),
		},
		&cli.StringFlag{
			Name:    "label.full_name",
			Usage:   "full name of the repository",
			Sources: cli.EnvVars("VELA_REPO_FULL_NAME"),
		},
		&cli.StringFlag{
			Name:    "label.url",
			Usage:   "direct url of the repository",
			Sources: cli.EnvVars("VELA_REPO_LINK"),
		},
		&cli.StringFlag{
			Name:    "label.build_link",
			Usage:   "direct Vela link to the build",
			Sources: cli.EnvVars("VELA_BUILD_LINK"),
		},
		&cli.StringFlag{
			Name:    "label.host",
			Usage:   "host that the image is built on",
			Sources: cli.EnvVars("VELA_BUILD_HOST"),
		},
		&cli.StringFlag{
			Name:    "label.custom",
			Usage:   "custom labels to add to the image in the form LABEL_NAME=ENV_KEY",
			Sources: cli.EnvVars("VELA_BUILD_CUSTOM_LABELS", "PARAMETER_CUSTOM_LABELS"),
		},
		&cli.StringSliceFlag{
			Name:    "label.topics",
			Usage:   "topics of the repository",
			Sources: cli.EnvVars("VELA_REPO_TOPICS"),
		},
	}

	err = app.Run(context.Background(), os.Args)
	if err != nil {
		logrus.Fatal(err)
	}
}

// run executes the plugin based off the configuration provided.
func run(ctx context.Context, c *cli.Command) error {
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

	// target type for build args
	var buildArgs []string

	argsStr := c.String("image.build_args")
	if len(argsStr) > 0 {
		buildArgsMap := make(map[string]string)

		// attempt to unmarshal to map
		err := json.Unmarshal([]byte(argsStr), &buildArgsMap)
		if err != nil {
			// fall back on splitting the string
			buildArgs = strings.Split(argsStr, ",")
		} else {
			// iterate through the build args map
			for key, value := range buildArgsMap {
				// add the build arg to the build args
				buildArgs = append(buildArgs, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	// target type for custom labels
	var customLabels []string

	labelsStr := c.String("label.custom")
	if len(labelsStr) > 0 {
		customLabelsMap := make(map[string]string)

		// attempt to unmarshal to map
		err := json.Unmarshal([]byte(labelsStr), &customLabelsMap)
		if err != nil {
			// fall back on splitting the string
			customLabels = strings.Split(labelsStr, ",")
		} else {
			// iterate through the custom labels map
			for key, value := range customLabelsMap {
				// add the custom label to the custom labels
				customLabels = append(customLabels, fmt.Sprintf("%s=%s", key, value))
			}
		}
	}

	// create the plugin
	p := &Plugin{
		// build configuration
		Build: &Build{
			Event:          c.String("build.event"),
			Sha:            c.String("build.sha"),
			SnapshotMode:   c.String("build.snapshot_mode"),
			Tag:            c.String("build.tag"),
			UseNewRun:      c.Bool("build.use_new_run"),
			TarPath:        c.String("build.tar_path"),
			SingleSnapshot: c.Bool("build.single_snapshot"),
			IgnoreVarRun:   c.Bool("build.ignore_var_run"),
			IgnorePath:     c.StringSlice("build.ignore_path"),
			LogTimestamp:   c.Bool("build.log_timestamps"),
		},
		// image configuration
		Image: &Image{
			Args:               buildArgs,
			Context:            c.String("image.context"),
			Dockerfile:         c.String("image.dockerfile"),
			Target:             c.String("image.target"),
			ForceBuildMetadata: c.Bool("image.force_build_metadata"),
			CustomPlatform:     c.String("image.custom_platform"),
		},
		// registry configuration
		Registry: &Registry{
			DryRun:             c.Bool("registry.dry_run"),
			Name:               c.String("registry.name"),
			Mirror:             c.String("registry.mirror"),
			Username:           c.String("registry.username"),
			Password:           c.String("registry.password"),
			PushRetry:          c.Int("registry.push_retry"),
			InsecureRegistries: c.StringSlice("registry.insecure_registries"),
			InsecurePull:       c.Bool("registry.insecure_pull"),
			InsecurePush:       c.Bool("registry.insecure_push"),
		},
		// repo configuration
		Repo: &Repo{
			AutoTag:           c.Bool("repo.auto_tag"),
			Cache:             c.Bool("repo.cache"),
			CacheName:         c.String("repo.cache_name"),
			Compression:       c.String("repo.compression"),
			CompressionLevel:  c.Int("repo.compression_level"),
			CompressedCaching: c.Bool("repo.compressed_caching"),
			Name:              c.String("repo.name"),
			Tags:              c.StringSlice("repo.tags"),
			TopicsFilter:      c.String("repo.topics_filter"),
			Label: &Label{
				AuthorEmail: c.String("label.author_email"),
				Commit:      c.String("label.commit"),
				Created:     time.Now().Format(time.RFC3339),
				FullName:    c.String("label.full_name"),
				Number:      c.Int("label.number"),
				Topics:      c.StringSlice("label.topics"),
				URL:         c.String("label.url"),
				BuildURL:    c.String("label.build_link"),
				Host:        c.String("label.host"),
				CustomSet:   customLabels,
			},
			Labels: c.StringSlice("repo.labels"),
		},
	}

	// check if repo auto tagging is enabled
	if p.Repo.AutoTag {
		p.Repo.ConfigureAutoTagBuildTags(p.Build)
	}

	// validate the plugin
	err := p.Validate()
	if err != nil {
		return err
	}

	// execute the plugin
	return p.Exec(ctx)
}
