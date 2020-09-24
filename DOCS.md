## Description

This plugin enables you to build and publish [Docker](https://www.docker.com/) images in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-kaniko

Registry: https://hub.docker.com/r/target/vela-kaniko

## Usage

_Notes:_

* The plugin supports reading all parameters via environment variables or files. Values set as a file take precedence over default values set from the environment.
* We do not recommended using latest for pipelines. Users should use pinned images to decrease volatility of external changes to their pipelines. 
* The [Snapshot mode](https://github.com/GoogleContainerTools/kaniko/releases/tag/v1.0.0) can help increase build times but it is recommend to follow Kanikos guidelines for picking the mode.

Sample of building and publishing an image:

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building an image without publishing:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
+     dry_run: true
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building and publishing an image with custom tags:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
+     tags:
+       - latest
+       - foobar
```

Sample of building and publishing an image with automatic tags:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
+     auto_tag: true
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building and publishing an image with build arguments:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
+     build_args:
+       - FOO=bar
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building and publishing an image with caching:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
+     cache: true
+     cache_repo: index.docker.io/octocat/hello-world
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building using a snapshot mode and publishing an image with caching:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
    parameters:
+     snapshot_mode: redo
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

## Secrets

**NOTE: Users should refrain from configuring sensitive information in your pipeline in plain text.**

You can use Vela secrets to substitute sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: true
+   secrets: [ docker_username, docker_password ]
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
-     username: octocat
-     password: superSecretPassword
```

## Parameters

**NOTE: Vela injects several variables, by default, that this plugin can load in automatically.**

The following parameters are used to configure the image:

| Name            | Description                                                        | Required | Default           |
| --------------- | ------------------------------------------------------------------ | -------- | ----------------- |
| `auto_tag`      | enables tagging of image automatically                             | `false`  | `false`           |
| `build_args`    | variables passed to image at build-time                            | `false`  | `N/A`             |
| `cache`         | enable caching of image layers                                     | `false`  | `false`           |
| `cache_repo`    | specific repo to enable caching for                                | `false`  | `N/A`             |
| `context`       | path to context for building the image                             | `true`   | `.`               |
| `dockerfile`    | path to the file for building the image                            | `true`   | `Dockerfile`      |
| `dry_run`       | enable building the image without publishing                       | `false`  | `false`           |
| `event`         | event generated for build                                          | `true`   | **set by Vela**   |
| `log_level`     | set the log level for the plugin                                   | `true`   | `info`            |
| `mirror`        | name of the mirror registry to use                                 | `false`  | `N/A`             |
| `password`      | password for communication with the registry                       | `true`   | `N/A`             |
| `registry`      | name of the registry for the repository                            | `true`   | `index.docker.io` |
| `repo`          | name of the repository for the image                               | `true`   | `N/A`             |
| `sha`           | SHA-1 hash generated for commit                                    | `true`   | **set by Vela**   |
| `snapshot_mode` | control how to snapshot the filesystem. - options (full|redo|time) | `false`  | **set by Vela**   |
| `tag`           | tag generated for build                                            | `false`  | **set by Vela**   |
| `tags`          | unique tags of the image                                           | `true`   | `latest`          |
| `username`      | user name for communication with the registry                      | `true`   | `N/A`             |

## Template

COMING SOON!

## Troubleshooting

Below are a list of common problems and how to solve them:
