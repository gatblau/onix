#
#  Onix Config Manager - Artisan Linux Container Image Builder Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM docker.io/centos:7

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL builder-for=CONTAINER-IMAGES
LABEL tools="buildah,artisan"

ARG APP_NAME=art

ENV UID=1000
ENV GID=1000

RUN yum update -y && \
    yum install buildah skopeo wget -y && \
    rm -rf /var/cache/yum && \
    yum clean all

COPY ./images/art-buildah/build.sh ./images/bin/output/$APP_NAME /usr/bin/

RUN ln -s /usr/bin/$APP_NAME /usr/bin/art && \
    chmod ug+x /usr/bin/build.sh && \
    chmod ug+x /usr/bin/$APP_NAME && \
    chmod ug+x /usr/bin/art && \
    mkdir -p /keys

CMD ["/usr/bin/build.sh"]