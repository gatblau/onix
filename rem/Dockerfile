#
#    Onix REMote Host - Remote Host Service for Pilot
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=remote

ENV UID=100000000
ENV GID=100000000

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    useradd -M -u $UID -g $GID $UNAME && \
    usermod --home / -u $UID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

WORKDIR /app

COPY ./bin/rem ./

USER $UNAME

CMD ["sh", "-c", "/app/rem"]

EXPOSE 8080/tcp
