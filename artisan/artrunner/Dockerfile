#
#    Onix Config Manager - Artisan Flow Runner Image
#    Copyright (c) 2018-Present by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=runner

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

COPY ./bin/linux/runner ./

USER $UNAME

CMD ["sh", "-c", "/app/runner"]

EXPOSE 8080/tcp
