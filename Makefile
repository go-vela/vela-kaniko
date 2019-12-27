# Copyright (c) 2019 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

build: binary-build

run: build docker-build docker-run

test: build docker-build docker-example

#################################
######      Go clean       ######
#################################

clean:

	@go mod tidy
	@go vet ./...
	@go fmt ./...
	@echo "I'm kind of the only name in clean energy right now"

#################################
######    Build Binary     ######
#################################

binary-build:

	GOOS=linux CGO_ENABLED=0 go build -o release/vela-docker github.com/go-vela/vela-docker

#################################
######    Docker Build     ######
#################################

docker-build:

	docker build --no-cache -t vela-docker:local .

#################################
######     Docker Run      ######
#################################

docker-run:

	docker run --rm \
		-e BUILD_EVENT \
		-e BUILD_COMMIT \
		-e BUILD_TAG \
		-e PARAMETER_REGISTRY \
		-e PARAMETER_REPO \
		-e PARAMETER_TAGS \
		-e PARAMETER_USERNAME \
		-e PARAMETER_PASSWORD \
		-e PARAMETER_DOCKERFILE \
		-e PARAMETER_DRY_RUN \
		-e PARAMETER_CONTEXT \
		-e PARAMETER_BUILD_ARGS \
		-e PARAMETER_LOG_LEVEL \
		-e PARAMETER_CACHE \
		-e PARAMETER_CACHE_REPO \
		-e PARAMETER_AUTO_TAG \
		-v $(pwd):/workspace
		vela-docker:local

docker-example:

	docker run --rm \
	-e BUILD_COMMIT=123abcdefg
	-e BUILD_EVENT=push
	-e PARAMETER_REGISTRY=index.docker.io
	-e PARAMETER_REPO=index.docker.io/target/vela-docker
	-e PARAMETER_TAGS=latest
	-e PARAMETER_USERNAME
	-e PARAMETER_PASSWORD
	-v $(pwd):/workspace
	vela-docker:local
