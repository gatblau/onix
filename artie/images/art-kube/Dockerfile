#
#  Onix Config Manager - Artie K8S Launcher Image
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM registry.access.redhat.com/ubi8/ubi-minimal

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL tools="kubectl, oc, artie"

ARG APP_NAME=artie
ARG UNAME=K8Sclient

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV ARTIE_HOME /usr/bin
ENV OKD_VERSION=4.5.0-0.okd-2020-10-15-235428

RUN microdnf update --disablerepo=* --enablerepo=ubi-8-appstream --enablerepo=ubi-8-baseos -y && \
    microdnf install shadow-utils.x86_64 -y && \
    microdnf install wget gzip tar -y && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    microdnf clean all

RUN wget https://github.com/openshift/okd/releases/download/$OKD_VERSION/openshift-client-linux-$OKD_VERSION.tar.gz && \
    tar -xf openshift-client-linux-$OKD_VERSION.tar.gz && \
    mv kubectl $ARTIE_HOME && \
    mv oc $ARTIE_HOME && \
    rm openshift-client-linux-$OKD_VERSION.tar.gz

COPY ./images/bin/output/$APP_NAME ./images/art-kube/buildconfig.yaml.tem $ARTIE_HOME/

RUN ln -s $ARTIE_HOME/$APP_NAME $ARTIE_HOME/artie && \
    chmod ug+x $ARTIE_HOME/$APP_NAME && \
    chmod ug+x $ARTIE_HOME/artie && \
    chmod ug+x $ARTIE_HOME/kubectl && \
    chmod ug+x $ARTIE_HOME/oc && \
    chmod a+r $ARTIE_HOME/buildconfig.yaml.tem

USER $UNAME