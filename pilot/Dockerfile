# Onix Config Manager - Pilot Sidecar
# Copyright (c) 2018-2020 by www.gatblau.org
# Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
# Contributors to this project, hereby assign copyright in this code to the project,
# to be licensed under the same terms as the rest of the code.

# This dockerfile encapsulates the build process for the Onix Web API
# The builder container is transient and downloads and install maven, package the Java app and extracts the
# Springboot uberjar files to improve startup times
# The release image copy the prepared app files from the builder image

# NOTE: see https://access.redhat.com/solutions/4643601
# dnf checks for subsciptions even if not needed (not using RHEL repos)

# compile stage: build Pilot
FROM golang as builder
WORKDIR /app
COPY . .
RUN go get . && CGO_ENABLED=0 GOOS=linux go build -o pilot .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
ENV UID 25
ENV GID 25
WORKDIR /app
COPY --from=builder /app/pilot /app/config.toml ./
USER UID:GID
CMD ["/app/pilot", "sidecar"]
