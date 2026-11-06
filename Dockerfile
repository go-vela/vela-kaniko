# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-kaniko:certs .    ##
#########################################################################

ARG KANIKO_IMAGE=target/kaniko-executor:debug-v1.25.19@sha256:8748e38b3bed2e7a91f13d07c1f6c64086a0430f511b59eadffa94876b3614b1

FROM alpine:3.24.4@sha256:93b41ef40b1677fce8fcc64ff4af09daddfacf5040afde618bc46c1e2d1cc7d4 AS certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-kaniko:local .    ##
##########################################################

FROM ${KANIKO_IMAGE}

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /workspace

COPY release/vela-kaniko /bin/vela-kaniko

ENTRYPOINT [ "/bin/vela-kaniko" ]
