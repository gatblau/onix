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

ARG UNAME=pilotctl

ENV UID=100000000
ENV GID=100000000

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    useradd -m -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all && \
    # keys volume mount
    mkdir -p /keys/ && \
    mkdir /app

WORKDIR /app

COPY ./bin/pilotctl .
COPY boot.sh .

RUN chown -R pilotctl:pilotctl /app

USER $UNAME

CMD ["sh", "-c", "bash boot.sh"]

EXPOSE 8080/tcp