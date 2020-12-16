#
#  Onix Config Manager - Artisan K8S Launcher Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL tools="kubectl,oc,artisan"

ARG APP_NAME=art
ARG UNAME=K8Sclient

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV ARTISAN_HOME /usr/bin
ENV OKD_VERSION=4.5.0-0.okd-2020-10-15-235428

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 -y && \
    microdnf install wget gzip tar -y && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

RUN wget https://github.com/openshift/okd/releases/download/$OKD_VERSION/openshift-client-linux-$OKD_VERSION.tar.gz && \
    tar -xf openshift-client-linux-$OKD_VERSION.tar.gz && \
    mv kubectl $ARTISAN_HOME && \
    mv oc $ARTISAN_HOME && \
    rm openshift-client-linux-$OKD_VERSION.tar.gz

COPY ./images/bin/output/$APP_NAME ${ARTISAN_HOME}
COPY ./images/art-kube/buildconfig.yaml.tem /tmp/

RUN ln -s $ARTISAN_HOME/$APP_NAME $ARTISAN_HOME/art && \
    chmod ug+x $ARTISAN_HOME/$APP_NAME && \
    chmod ug+x $ARTISAN_HOME/art && \
    chmod ug+x $ARTISAN_HOME/kubectl && \
    chmod ug+x $ARTISAN_HOME/oc && \
    chmod a+xwr /tmp/buildconfig.yaml.tem

USER $UNAME