# Onix Config Manager - Dbman
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

# building stage: compiles dbman
FROM golang as builder
ARG DB_PLUGIN_PREFIX="dbman-db-"
WORKDIR /app
COPY . .
RUN cd ./plugin/pgsql && \
    export GOOS=linux; export GOARCH=amd64; go build -o ${DB_PLUGIN_PREFIX}pgsql && \
    cd ../.. && \
    mv ./plugin/pgsql/${DB_PLUGIN_PREFIX}pgsql ./${DB_PLUGIN_PREFIX}pgsql && \
    export GOOS=linux; export GOARCH=amd64; go build

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi-minimal
LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
ARG UNAME=dbman
ENV UID=1000
ENV GID=1000
RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all
WORKDIR /app
# copy dbman binaries
COPY --from=builder /app/dbman /app/dbman-db-* /app/
# copy config files to user home
# 1. if running in podman or docker user home is "/home/dbman/"
COPY --from=builder /app/.dbman.toml /app/.dbman_default.toml /home/dbman/
# 2. if running in openshift user home is "/"
COPY --from=builder /app/.dbman.toml /app/.dbman_default.toml /
USER $UNAME
CMD ["/app/dbman", "serve"]
EXPOSE 8085/tcp
