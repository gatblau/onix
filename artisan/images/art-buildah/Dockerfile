#
#  Onix Config Manager - Artisan Linux Container Image Builder Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM docker.io/centos:8.3.2011

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL tools="buildah,skopeo,artisan"

RUN dnf update -y && \
    dnf install buildah skopeo -y && \
    rm -rf /var/cache/dnf && \
    dnf clean all

COPY art /usr/bin/art

RUN chmod ug+x /usr/bin/art && mkdir -p /keys
