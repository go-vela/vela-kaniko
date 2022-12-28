## Description

This plugin enables you to build and publish [Docker](https://www.docker.com/) images in a Vela pipeline.

Source Code: https://github.com/go-vela/vela-kaniko

Registry: https://hub.docker.com/r/target/vela-kaniko

## Usage

> **NOTE:**
>
> Users should refrain from using latest as the tag for the Docker image.
>
> It is recommended to use a semantically versioned tag instead.

Sample of building and publishing an image:

```yaml
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building an image without publishing:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
    parameters:
+     dry_run: true
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of attempting the publishing of an image three times:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
    parameters:
+     push_retry: 3
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building and publishing an image with custom tags:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
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
    pull: always
    parameters:
+     auto_tag: true
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Depending on the type of event, the image will be tagged as follows:

* tag event (using `v1.0.0` as an example):
    * `index.docker.io/octocat/hello-world:latest`
    * `index.docker.io/octocat/hello-world:v1.0.0`

* all other events:
    * `index.docker.io/octocat/hello-world:latest`
    * `index.docker.io/octocat/hello-world:eeea105fed7fc11bda4b43a00edfc49a5c982968`


Sample of building and publishing an image with build arguments:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
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
    pull: always
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
    pull: always
    parameters:
+     snapshot_mode: redo
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Sample of building using a custom platform:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
+     custom_platform: linux/arm64/v8
```

> **NOTE:** This option will only work if your Vela worker is configured appropriately.

## Secrets

> **NOTE:** Users should refrain from configuring sensitive information in your pipeline in plain text.

### Internal

Users can use [Vela internal secrets](https://go-vela.github.io/docs/tour/secrets/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
+   secrets: [ kaniko_username, kaniko_password ]
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
-     username: octocat
-     password: superSecretPassword
```

> This example will add the secrets to the `publish_hello-world` step as environment variables:
>
> * `KANIKO_USERNAME=<value>`
> * `KANIKO_PASSWORD=<value>`

### External

The plugin accepts the following files for authentication:

| Parameter  | Volume Configuration                                                |
| ---------- | ------------------------------------------------------------------- |
| `password` | `/vela/parameters/kaniko/password`, `/vela/secrets/kaniko/password` |
| `username` | `/vela/parameters/kaniko/username`, `/vela/secrets/kaniko/username` |

Users can use [Vela external secrets](https://go-vela.github.io/docs/concepts/pipeline/secrets/origin/) to substitute these sensitive values at runtime:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
+   secrets: [ kaniko_username, kaniko_password ]
    parameters:
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
-     username: octocat
-     password: superSecretPassword
```

> This example will read the secret values in the volume stored at `/vela/secrets/`

## Parameters

> **NOTE:**
>
> The plugin supports reading all parameters via environment variables or files.
>
> Any values set from a file take precedence over values set from the environment.
>
> The [Snapshot mode](https://github.com/GoogleContainerTools/kaniko/releases/tag/v1.0.0) can help improve performance but it is recommend to follow Kaniko's guidelines for picking the mode.

The following parameters are used to configure the image:

| Name                   | Description                                                        | Required | Default           | Environment Variables                                             |
| ---------------------- | ------------------------------------------------------------------ | -------- | ----------------- | ----------------------------------------------------------------- |
| `auto_tag`             | enables automatic tagging of images (tag or sha, and `latest`)     | `false`  | `false`           | `PARAMETER_AUTO_TAG`<br>`KANIKO_AUTO_TAG`                         |
| `build_args`           | variables passed to image at build-time                            | `false`  | `N/A`             | `PARAMETER_BUILD_ARGS`<br>`KANIKO_BUILD_ARGS`                     |
| `cache`                | enable caching of image layers                                     | `false`  | `false`           | `PARAMETER_CACHE`<br>`KANIKO_CACHE`                               |
| `cache_repo`           | specific repo to enable caching for                                | `false`  | `N/A`             | `PARAMETER_CACHE_REPO`<br>`KANIKO_CACHE_REPO`                     |
| `context`              | path to context for building the image                             | `true`   | `.`               | `PARAMETER_CONTEXT`<br>`KANIKO_CONTEXT`                           |
| `dockerfile`           | path to the file for building the image                            | `true`   | `Dockerfile`      | `PARAMETER_DOCKERFILE`<br>`KANIKO_DOCKERFILE`                     |
| `dry_run`              | enable building the image without publishing                       | `false`  | `false`           | `PARAMETER_DRY_RUN`<br>`KANIKO_DRY_RUN`                           |
| `event`                | event generated for build                                          | `true`   | **set by Vela**   | `PARAMETER_EVENT`<br>`KANIKO_EVENT`<br>`VELA_BUILD_EVENT`         |
| `force_build_metadata` | enable force adding metadata layers to build image                 | `false`  | `false`           | `PARAMETER_FORCE_BUILD_METADATA`<br>`KANIKO_FORCE_BUILD_METADATA` |
| `labels`               | unique labels to add to the image                                  | `false`  | `N/A`             | `PARAMETER_LABELS`<br>`KANIKO_LABELS`                             |
| `log_level`            | set the log level for the plugin                                   | `true`   | `info`            | `PARAMETER_LOG_LEVEL`<br>`KANIKO_LOG_LEVEL`                       |
| `mirror`               | name of the mirror registry to use                                 | `false`  | `N/A`             | `PARAMETER_MIRROR`<br>`KANIKO_MIRROR`                             |
| `password`             | password for communication with the registry                       | `true`   | `N/A`             | `PARAMETER_PASSWORD`<br>`KANIKO_PASSWORD`<br>`DOCKER_PASSWORD`    |
| `push_retry`           | number of retries for pushing an image to a remote destination     | `false`  | `0`               | `PARAMETER_PUSH_RETRY`<br>`KANIKO_PUSH_RETRY`                     |
| `registry`             | name of the registry for the repository                            | `true`   | `index.docker.io` | `PARAMETER_REGISTRY`<br>`KANIKO_REGISTRY`                         |
| `repo`                 | name of the repository for the image                               | `true`   | `N/A`             | `PARAMETER_REPO`<br>`KANIKO_REPO`                                 |
| `sha`                  | SHA-1 hash generated for commit                                    | `true`   | **set by Vela**   | `PARAMETER_SHA`<br>`KANIKO_SHA`<br>`VELA_BUILD_COMMIT`            |
| `snapshot_mode`        | control how to snapshot the filesystem. - options (full|redo|time) | `false`  | `N/A`             | `PARAMETER_SNAPSHOT_MODE`<br>`KANIKO_SNAPSHOT_MODE`               |
| `use_new_run`          | use experimental run implementation for detecting changes without requiring file system snapshots | `false` | `N/A` | `PARAMETER_USE_NEW_RUN`<br>`KANIKO_USE_NEW_RUN` |
| `tag`                  | tag generated for build                                            | `false`  | **set by Vela**   | `PARAMETER_TAG`<br>`KANIKO_TAG`<br>`VELA_BUILD_TAG`               |
| `tags`                 | unique tags of the image                                           | `true`   | `latest`          | `PARAMETER_TAGS`<br>`KANIKO_TAGS`                                 |
| `target`               | set the target build stage for the image                           | `false`  | `N/A`             | `PARAMETER_TARGET`<br>`KANIKO_TARGET`                             |
| `username`             | user name for communication with the registry                      | `true`   | `N/A`             | `PARAMETER_USERNAME`<br>`KANIKO_USERNAME`<br>`DOCKER_USERNAME`    |
| `custom_platform`      | set the custom platform for the image                               | `false`  | `N/A`             | `PARAMETER_CUSTOM_PLATFORM`<br>`KANIKO_CUSTOM_PLATFORM`           |
| `insecure_registries`  | insecure docker registries to push or pull to/from                  | `false`  | `empty slice`     | `PARAMETER_INSECURE_REGISTRIES`<br>`KANIKO_INSECURE_REGISTRIES`   |
| `insecure_pull`        | enable pulling from any insecure registry                           | `false`  | `false`           | `PARAMETER_INSECURE_PULL`<br>`KANIKO_INSECURE_PULL`               |
| `insecure_push`        | enable pushing to any insecure registry                             | `false`  | `false`           | `PARAMETER_INSECURE_PUSH`<br>`KANIKO_INSECURE_PUSH`               |

## Template

COMING SOON!

## Troubleshooting

You can start troubleshooting this plugin by tuning the level of logs being displayed:

```diff
steps:
  - name: publish_hello-world
    image: target/vela-kaniko:latest
    pull: always
    parameters:
+     log_level: trace
      registry: index.docker.io
      repo: index.docker.io/octocat/hello-world
```

Below are a list of common problems and how to solve them:
