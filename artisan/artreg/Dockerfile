#
#    Onix Config Manager - Artisan Package Registry Service Image
#    Copyright (c) 2018-Present by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=artisan

ENV UID=1000
ENV GID=1000

# the location of the artisan registry files in the image
# NB if changed, ensure there is a trailing slash at the end of the path
ENV ARTISAN_HOME=/opt/

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 && \
    groupadd -g $GID -o $UNAME && \
    # -M create the user with no /home
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all && \
    # create new Artisan specific directory and sub folders to support pessimistic locking for metadata changes within a
    # virtual repository
    mkdir -p \
        ${ARTISAN_HOME}.artisan/locks \
        ${ARTISAN_HOME}.artisan/tmp \
        ${ARTISAN_HOME}.artisan/hooks && \
    chown -R ${UNAME} ${ARTISAN_HOME}.artisan

WORKDIR /app

COPY ./bin/linux/artreg ./

USER $UNAME

CMD ["sh", "-c", "/app/artreg"]

EXPOSE 8082/tcp
