#
#    Onix Config Manager - Doorman Proxy
#    Copyright (c) 2018-Present by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=dproxy

ENV UID=1000
ENV GID=1000
ENV OX_HTTP_PORT=9998

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    # -M create the user with no /home
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

WORKDIR /app

COPY ./bin/linux/dproxy ./

USER $UNAME

CMD ["sh", "-c", "/app/dproxy"]

EXPOSE 9998/tcp
