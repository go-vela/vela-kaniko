# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-kaniko:certs .    ##
#########################################################################

FROM alpine:3.23.3@sha256:25109184c71bdad752c8312a8623239686a9a2071e8825f20acb8f2198c3f659 as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-kaniko:local .    ##
##########################################################

# Allow the kaniko base image to be overridden via build arg
# renovate: datasource=github-releases depName=chainguard-forks/kaniko
ARG KANIKO_IMAGE=target/kaniko-executor:debug-v1.24.0

FROM ${KANIKO_IMAGE}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /workspace

COPY release/vela-kaniko /bin/vela-kaniko

ENTRYPOINT [ "/bin/vela-kaniko" ]
