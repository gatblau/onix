#
#  Onix Config Manager - Artie Builder for Java 11
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM centos:7

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL builder-for=JAVA-11
LABEL tools="maven,jdk11"

ARG APP_NAME=artie
ARG UNAME=builder

ENV UID=1000
ENV GID=1000
ENV APP=$APP_NAME
ENV MAVEN_VERSION 3.6.3
ENV JAVA_HOME /usr/lib/jvm/jre-openjdk
ENV MAVEN_HOME /usr/lib/mvn
ENV ARTIE_HOME /usr/bin
ENV PATH $MAVEN_HOME/bin:$PATH

RUN yum update -y && \
    yum install -y java-11-openjdk-devel && \
    yum install -y wget && \
    yum install -y git unzip && \
    groupadd -g $GID -o $UNAME && useradd -M -u $UID -g $GID $UNAME && \
    rm -rf /var/cache/yum && \
    yum clean all

RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

COPY ./images/bin/output/$APP_NAME ./images/java-builder/start.sh $ARTIE_HOME/

RUN ln -s $ARTIE_HOME/$APP_NAME /usr/bin/artie && \
    chmod ug+x $ARTIE_HOME/$APP_NAME && \
    chmod ug+x $ARTIE_HOME/artie && \
    chmod ug+x $ARTIE_HOME/start.sh && \
    # folder below to mount keys in container
    mkdir -p /keys

USER $UNAME

ENTRYPOINT ["sh", "start.sh"]