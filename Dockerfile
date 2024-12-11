# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-kaniko:certs .    ##
#########################################################################

FROM alpine:3.21.0@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45 as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-kaniko:local .    ##
##########################################################

FROM gcr.io/kaniko-project/executor:v1.23.2-debug@sha256:c3109d5926a997b100c4343944e06c6b30a6804b2f9abe0994d3de6ef92b028e

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /workspace

COPY release/vela-kaniko /bin/vela-kaniko

ENTRYPOINT [ "/bin/vela-kaniko" ]
