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
WORKDIR /app
COPY . .
RUN go get . && CGO_ENABLED=0 GOOS=linux go build -o dbman .

# package stage: copy the binary into the deployment image
FROM registry.access.redhat.com/ubi8/ubi
LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
WORKDIR /app
COPY --from=builder /app/dbman ./
COPY --from=builder /app/.dbman.toml /app/.dbman_default.toml /
#RUN rpm -Uvh https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm && microdnf update -y && rm -rf /var/cache/yum && microdnf install postgresql12 && microdnf clean all
USER 20
CMD ["/app/dbman", "serve"]
EXPOSE 8085/tcp
