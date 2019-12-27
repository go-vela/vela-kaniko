# Copyright (c) 2019 Target Brands, Inc. All rights reserved.
#
# Use of this source code is governed by the LICENSE file in this repository.

FROM gcr.io/kaniko-project/executor:debug-v0.13.0

WORKDIR /workspace

COPY release/vela-docker /bin/vela-docker

ENTRYPOINT [ "/bin/vela-docker" ]
