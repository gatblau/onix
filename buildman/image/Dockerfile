#
#  Onix Config Manager - Build Manager
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM centos:7

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"

ARG APP_NAME=buildman
ARG UNAME=buildman

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV APP_HOME /app

RUN yum update -y && \
    yum install shadow-utils.x86_64 -y && \
    yum install skopeo -y && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    yum clean all && \
    mkdir /conf

COPY ./image/bin/output/$APP_NAME $APP_HOME/
COPY ./image/policy.json /conf/

RUN chmod ug+x $APP_HOME/$APP_NAME && chmod ug+r /conf/policy.json

USER $UNAME

CMD ["sh", "-c", "${APP_HOME}/${APP}"]