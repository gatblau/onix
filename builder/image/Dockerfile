#
#  Onix Config Manager - Onix Build Manager Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM centos:8

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG APP_NAME=buildman
ARG UNAME=builder

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV APP_HOME /app

RUN dnf update && \
    dnf install shadow-utils.x86_64 -y && \
    dnf install skopeo -y && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    dnf clean all

COPY ./image/bin/output/$APP_NAME $APP_HOME/

RUN chmod ug+x $APP_HOME/$APP_NAME

USER $UNAME

CMD ["sh", "-c", "${APP_NAME}/${APP}"]