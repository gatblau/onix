#
#  Onix Config Manager - Artie Linux Container Image Builder Image
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
LABEL tools="buildah"

ARG APP_NAME=artie
#ARG UNAME=builder

ENV UID=1000
ENV GID=1000

RUN yum update -y && \
    yum install buildah wget -y && \
#    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    yum clean all

COPY ./images/art-buildah/build.sh ./images/bin/output/$APP_NAME /usr/bin/

RUN ln -s /usr/bin/$APP_NAME /usr/bin/artie && \
    chmod ug+x /usr/bin/build.sh && \
    chmod ug+x /usr/bin/$APP_NAME && \
    chmod ug+x /usr/bin/artie && \
    # folder below to mount keys in container
    mkdir -p /keys

#USER $UNAME

CMD ["/usr/bin/build.sh"]