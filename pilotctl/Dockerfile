#
#    Onix Pilot Host Control Service
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL artisan.svc.manifest="/app/svc.yaml"

ARG UNAME=pilotctl

ENV UID=100
ENV GID=100

ENV SYNC_PATH=/sync

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all && \
    mkdir -p ${SYNC_PATH} && \
    chown -R ${UNAME} ${SYNC_PATH}

USER $UNAME

WORKDIR /app

COPY ./bin/pilotctl ./svc.yaml ./

# mount for PGP signing key
VOLUME /keys

# mount for event receivers configuration file
VOLUME /conf

# mount for sync files
VOLUME /sync

CMD ["sh", "-c", "/app/pilotctl"]

EXPOSE 8080/tcp