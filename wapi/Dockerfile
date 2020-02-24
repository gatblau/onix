#
#    Onix Config Manager - Copyright (c) 2018-2019 by www.gatblau.org
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#    Unless required by applicable law or agreed to in writing, software distributed under
#    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#    either express or implied.
#    See the License for the specific language governing permissions and limitations under the License.
#
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
# the builder image
FROM openjdk:15-jdk-alpine as builder
RUN apk add unzip && rm -rf /var/cache/apk/*
ENV MAVEN_VERSION 3.6.1
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION /usr/lib/mvn
WORKDIR /app
COPY . .
RUN rm -rf ./target && mvn -f pom.xml package && unzip -o ./target/*.jar -d /tmp
# the final deployment image
FROM quay.io/gatblau/openjdk:13-j9-ubi8-min
MAINTAINER Gatblau <onix@gatblau.org>
LABEL author="gatblau.org"
COPY --from=builder /tmp/BOOT-INF/lib /app/lib
COPY --from=builder /tmp/META-INF /app/META-INF
COPY --from=builder /tmp/BOOT-INF/classes /app
USER 20
ENTRYPOINT ["java","-cp","app:app/lib/*","org/gatblau/onix/App"]
