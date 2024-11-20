// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-vela/vela-kaniko/version"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	_ "github.com/joho/godotenv/autoload"
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
	app := cli.NewApp()

	// Plugin Information

	app.Name = "vela-kaniko"
	app.HelpName = "vela-kaniko"
	app.Usage = "Vela Kaniko plugin for building and publishing images"
	app.Copyright = "Copyright 2019 Target Brands, Inc. All rights reserved."
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
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_USE_NEW_RUN", "KANIKO_USE_NEW_RUN", "VELA_BUILD_USE_NEW_RUN"},
			FilePath: "/vela/parameters/kaniko/use_new_run,/vela/secrets/kaniko/use_new_run",
			Name:     "build.use_new_run",
			Usage:    "use the experimental run implementation for detecting changes without requiring file system snapshots. In some cases, this may improve build performance by 75%.",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_TAR_PATH", "KANIKO_TAR_PATH", "VELA_BUILD_TAR_PATH"},
			FilePath: "/vela/parameters/kaniko/tar_path,/vela/secrets/kaniko/tar_path",
			Name:     "build.tar_path",
			Usage:    "If set, the image will be saved as a tarball at that path. ",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_SINGLE_SNAPSHOT", "KANIKO_SINGLE_SNAPSHOT", "VELA_BUILD_SINGLE_SNAPSHOT"},
			FilePath: "/vela/parameters/kaniko/single_snapshot,/vela/secrets/kaniko/single_snapshot",
			Name:     "build.single_snapshot",
			Usage:    "takes a single snapshot of the filesystem at the end of the build, so only one layer will be appended to the base image",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_IGNORE_VAR_RUN", "KANIKO_IGNORE_VAR_RUN", "VELA_IGNORE_VAR_RUN"},
			FilePath: "/vela/parameters/kaniko/ignore_var_run,/vela/secrets/kaniko/ignore_var_run",
			Name:     "build.ignore_var_run",
			Usage:    "By default Kaniko ignores /var/run when taking image snapshot. Include this parameter to preserve /var/run/* in destination image.",
			Value:    true,
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_IGNORE_PATH", "KANIKO_IGNORE_PATH"},
			FilePath: "/vela/parameters/kaniko/ignore_path,/vela/secrets/kaniko/ignore_path",
			Name:     "build.ignore_path",
			Usage:    "ignore paths when taking image snapshot",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_LOG_TIMESTAMPS", "KANIKO_LOG_TIMESTAMPS"},
			FilePath: "/vela/parameters/kaniko/log_timestamps,/vela/secrets/kaniko/log_timestamps",
			Name:     "build.log_timestamps",
			Usage:    "enable timestamps in logs",
		},

		// Image Flags

		&cli.StringFlag{
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
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_FORCE_BUILD_METADATA", "KANIKO_FORCE_BUILD_METADATA"},
			FilePath: "/vela/parameters/kaniko/force_build_metadata,/vela/secrets/kaniko/force_build_metadata",
			Name:     "image.force_build_metadata",
			Usage:    "enables force adding metadata layers to build image",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_CUSTOM_PLATFORM", "KANIKO_CUSTOM_PLATFORM"},
			FilePath: "/vela/parameters/kaniko/custom_platform,/vela/secrets/kaniko/custom_platform",
			Name:     "image.custom_platform",
			Usage:    "custom platform for the image",
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
			FilePath: "/vela/parameters/kaniko/username,/vela/secrets/kaniko/username,/vela/secrets/managed-auth/username",
			Name:     "registry.username",
			Usage:    "user name for communication with the registry",
		},
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_PASSWORD", "KANIKO_PASSWORD", "DOCKER_PASSWORD"},
			FilePath: "/vela/parameters/kaniko/password,/vela/secrets/kaniko/password,/vela/secrets/managed-auth/password",
			Name:     "registry.password",
			Usage:    "password for communication with the registry",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_PUSH_RETRY", "KANIKO_PUSH_RETRY"},
			FilePath: "/vela/parameters/kaniko/push_retry,/vela/secrets/kaniko/push_retry",
			Name:     "registry.push_retry",
			Usage:    "number of retries for pushing an image to a remote destination",
		},
		&cli.StringSliceFlag{
			EnvVars:  []string{"PARAMETER_INSECURE_REGISTRIES", "KANIKO_INSECURE_REGISTRIES"},
			FilePath: "/vela/parameters/kaniko/insecure_registries,/vela/secrets/kaniko/insecure_registries",
			Name:     "registry.insecure_registries",
			Usage:    "insecure registries to push & pull from",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_INSECURE_PULL", "KANIKO_INSECURE_PULL"},
			FilePath: "/vela/parameters/kaniko/insecure_pull,/vela/secrets/kaniko/insecure_pull",
			Name:     "registry.insecure_pull",
			Usage:    "enable pulling from insecure registries",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_INSECURE_PUSH", "KANIKO_INSECURE_PUSH"},
			FilePath: "/vela/parameters/kaniko/insecure_push,/vela/secrets/kaniko/insecure_push",
			Name:     "registry.insecure_push",
			Usage:    "enable pushing to insecure registries",
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
			EnvVars:  []string{"PARAMETER_COMPRESSION", "KANIKO_COMPRESSION"},
			FilePath: "/vela/parameters/kaniko/compression,/vela/secrets/kaniko/compression",
			Name:     "repo.compression",
			Usage:    "set the compression type - gzip (default) or zstd",
		},
		&cli.IntFlag{
			EnvVars:  []string{"PARAMETER_COMPRESSION_LEVEL", "KANIKO_COMPRESSION_LEVEL"},
			FilePath: "/vela/parameters/kaniko/compression_level,/vela/secrets/kaniko/compression_level",
			Name:     "repo.compression_level",
			Usage:    "set the compression level (1-9, inclusive)",
		},
		&cli.BoolFlag{
			EnvVars:  []string{"PARAMETER_COMPRESSED_CACHING", "KANIKO_COMPRESSED_CACHING"},
			FilePath: "/vela/parameters/kaniko/compressed_caching,/vela/secrets/kaniko/compressed_caching",
			Name:     "repo.compressed_caching",
			Usage:    "when set to false, will prevent tar compression for cached layers",
			Value:    true,
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
		&cli.StringFlag{
			EnvVars:  []string{"PARAMETER_REPO_TOPICS_FILTER", "KANIKO_REPO_TOPICS_FILTER"},
			FilePath: "/vela/parameters/kaniko/repo_topics_filter,/vela/secrets/kaniko/repo_topics_filter",
			Name:     "repo.topics_filter",
			Usage:    "filter to restrict which repository topics to include in label",
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
		&cli.StringFlag{
			EnvVars: []string{"VELA_BUILD_LINK"},
			Name:    "label.build_link",
			Usage:   "direct Vela link to the build",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_BUILD_HOST"},
			Name:    "label.host",
			Usage:   "host that the image is built on",
		},
		&cli.StringFlag{
			EnvVars: []string{"VELA_BUILD_CUSTOM_LABELS", "PARAMETER_CUSTOM_LABELS"},
			Name:    "label.custom",
			Usage:   "custom labels to add to the image in the form LABEL_NAME=ENV_KEY",
		},
		&cli.StringSliceFlag{
			EnvVars: []string{"VELA_REPO_TOPICS"},
			Name:    "label.topics",
			Usage:   "topics of the repository",
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
	return p.Exec()
}
