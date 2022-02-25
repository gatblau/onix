#
#    Onix Config Manager - Artisan's Doorman
#    Copyright (c) 2018-Present by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG UNAME=doorman

ENV UID=1000
ENV GID=1000

# the location of the artisan registry files in the image
# NB if changed, ensure there is a trailing slash at the end of the path
ENV ARTISAN_HOME=/opt/
ENV OX_HTTP_PORT=9999

RUN dnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    dnf install -y  https://dl.fedoraproject.org/pub/epel/epel-release-latest-8.noarch.rpm && \
    dnf install -y clamav clamd clamav-update && \
    sed -i -e "s/^Example/#Example/" /etc/clamd.d/scan.conf && \
    sed -i -e "s/^Example/#Example/" /etc/freshclam.conf && \
    sed -i 's/#LocalSocket \/run/LocalSocket \/run/g' /etc/clamd.d/scan.conf && \
    groupadd -g $GID -o $UNAME && \
    # -M create the user with no /home
    useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    dnf clean all && \
    mkdir -p \
        ${ARTISAN_HOME}.artisan/locks \
        ${ARTISAN_HOME}.artisan/tmp \
        ${ARTISAN_HOME}.artisan/hooks && \
    chown -R ${UNAME} ${ARTISAN_HOME}.artisan

WORKDIR /app

COPY ./bin/linux/doorman ./

# clam needs root container
#USER $UNAME

CMD ["sh", "-c", "/app/doorman"]

EXPOSE 9999/tcp
