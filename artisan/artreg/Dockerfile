#
#    Onix Config Manager - Artisan Package Registry Service Image
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=artisan

ENV UID=100000000
ENV GID=100000000

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    # user home under /home/artisan
    useradd -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all && \
    # pessimistic locking for metadata changes within a virtual repository
    mkdir -p /home/artisan/locks && chown $UNAME /home/artisan/locks && chmod ug+wr /home/artisan/locks && \
    # temprorary file operations folder
    mkdir -p /home/artisan/tmp && chown $UNAME /home/artisan/tmp && chmod ug+wr /home/artisan/tmp && \
    # web hooks configuration
    mkdir -p /home/artisan/hooks && chown $UNAME /home/artisan/hooks && chmod ug+wr /home/artisan/hooks

WORKDIR /app

COPY ./bin/linux/artreg ./

USER $UNAME

CMD ["sh", "-c", "/app/artreg"]

EXPOSE 8082/tcp
