# SPDX-License-Identifier: Apache-2.0

#########################################################################
##    docker build --no-cache --target certs -t vela-kaniko:certs .    ##
#########################################################################

FROM alpine:3.20.3@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d as certs

RUN apk add --update --no-cache ca-certificates

##########################################################
##    docker build --no-cache -t vela-kaniko:local .    ##
##########################################################

FROM gcr.io/kaniko-project/executor:v1.23.2-debug@sha256:c3109d5926a997b100c4343944e06c6b30a6804b2f9abe0994d3de6ef92b028e

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /workspace

COPY release/vela-kaniko /bin/vela-kaniko

ENTRYPOINT [ "/bin/vela-kaniko" ]
