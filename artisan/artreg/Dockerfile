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
    # -M create the user with no /home
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all && \
    # pessimistic locking for metadata changes within a virtual repository
    mkdir -p /.artisan/locks && chown $UNAME /.artisan/locks && chmod ug+wr /.artisan/locks && \
    # temprorary file operations folder
    mkdir -p /.artisan/tmp && chown $UNAME /.artisan/tmp && chmod ug+wr /.artisan/tmp && \
    # web hooks configuration
    mkdir -p /.artisan/hooks && chown $UNAME /.artisan/hooks && chmod ug+wr /.artisan/hooks

WORKDIR /app

COPY ./bin/linux/artreg ./

USER $UNAME

CMD ["sh", "-c", "/app/artreg"]

EXPOSE 8082/tcp
