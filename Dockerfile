# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-kaniko:certs .    ##
#########################################################################

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c as certs

RUN apk add --update --no-cache ca-certificates

#########################################################################
##    Build Kaniko executor from Chainguard source                     ##
#########################################################################

FROM golang:1.24 as kaniko-builder

WORKDIR /go/src/github.com/chainguard-dev/kaniko

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG KANIKO_VERSION=v1.25.0
RUN git clone --depth 1 --branch ${KANIKO_VERSION} https://github.com/chainguard-dev/kaniko.git .

ENV CGO_ENABLED=0
ENV GOBIN=/usr/local/bin

RUN CGO_ENABLED=0 make out/executor

##########################################################
##    docker build --no-cache -t vela-kaniko:local .    ##
##########################################################

FROM alpine:3.21.3@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=kaniko-builder /go/src/github.com/chainguard-dev/kaniko/out/executor /kaniko/executor

RUN mkdir -p /workspace /kaniko/.docker /kaniko/ssl/certs && \
    chmod 777 /workspace /kaniko /kaniko/.docker

ENV HOME /root
ENV USER root
ENV PATH /usr/local/bin:/kaniko
ENV SSL_CERT_DIR /kaniko/ssl/certs
ENV DOCKER_CONFIG /kaniko/.docker

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /kaniko/ssl/certs/ca-certificates.crt

WORKDIR /workspace

COPY release/vela-kaniko /bin/vela-kaniko

ENTRYPOINT [ "/bin/vela-kaniko" ]
