#
#  Onix Config Manager - Artisan Node Application Builder Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM quay.io/gatblau/node:14.15.1

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL builder-for=NODE-14
LABEL tools="node,artisan"

ARG APP_NAME=art
ARG UNAME=builder

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV ARTISAN_HOME=/usr/bin

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

COPY ./images/bin/output/$APP_NAME ./images/start.sh $ARTISAN_HOME/

RUN ln -s $ARTISAN_HOME/$APP_NAME $ARTISAN_HOME/art && \
    chmod ug+x $ARTISAN_HOME/$APP_NAME && \
    chmod ug+x $ARTISAN_HOME/art && \
    chmod ug+x $ARTISAN_HOME/start.sh && \
    # folder below to mount keys in container
    mkdir -p /keys

USER $UNAME

ENTRYPOINT ["sh", "start.sh"]
