#
#  Onix Config Manager - Artie Builder for Java 11
#  Copyright (c) 2018-2020 by www.gatblau.org
#  Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#
#  Contributors to this project, hereby assign copyright in this code to the project,
#  to be licensed under the same terms as the rest of the code.
#
FROM centos:8

LABEL author="gatblau"
LABEL maintainer="onix@gatblau.org"
LABEL builder-for=JAVA-11
LABEL tools="maven,jdk11"

ENV MAVEN_VERSION 3.6.3
ENV JAVA_HOME /usr/lib/jvm/jre-openjdk
ENV MAVEN_HOME /usr/lib/mvn
ENV ARTIE_HOME /usr/bin
ENV PATH $MAVEN_HOME/bin:$PATH

RUN dnf install -y java-11-openjdk-devel && \
    dnf install -y wget

RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
    mv apache-maven-$MAVEN_VERSION $MAVEN_HOME

COPY artie start.sh $ARTIE_HOME/

RUN chmod 777 $ARTIE_HOME/artie && chmod 777 $ARTIE_HOME/start.sh

CMD ["sh", "start.sh"]