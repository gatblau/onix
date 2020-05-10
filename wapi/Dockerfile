#
#    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org
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
# This dockerfile encapsulates the build process for the Onix Web API
# The builder container is transient and downloads and install maven, package the Java app and extracts the
# Springboot uberjar files to improve startup times
# The release image copy the prepared app files from the builder image

# the builder transient container
FROM openjdk:15-jdk-alpine as builder
RUN apk add unzip && rm -rf /var/cache/apk/*
ENV MAVEN_VERSION 3.6.3
ENV MAVEN_HOME /usr/lib/mvn
ENV PATH $MAVEN_HOME/bin:$PATH
# download and install maven in the build container
RUN wget http://archive.apache.org/dist/maven/maven-3/$MAVEN_VERSION/binaries/apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  tar -zxvf apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  rm apache-maven-$MAVEN_VERSION-bin.tar.gz && \
  mv apache-maven-$MAVEN_VERSION /usr/lib/mvn
# define a working folder within the build container
WORKDIR /app
# copy the java project into the /app folder
COPY . .
# 1. package the app skipping the cucumber integration tests
#   (as it requires connectovity to database, that is not available within this process)
# 2. unzip the sprinboot uberjar into a /tmp folder ready to be copied by the release image
RUN mvn -Dmaven.test.skip=true -f pom.xml package && unzip -o ./target/*.jar -d /tmp

# the final release image
# uses universal base image rhel 8 minimal with OpenJ9 JVM and Open JDK 14
FROM quay.io/gatblau/openjdk:14-j9-ubi8-min
MAINTAINER Gatblau <onix@gatblau.org>
LABEL author="gatblau.org"
# copy the unzipped application files to the
COPY --from=builder /tmp/BOOT-INF/lib /app/lib
COPY --from=builder /tmp/META-INF /app/META-INF
COPY --from=builder /tmp/BOOT-INF/classes /app
USER 20
ENTRYPOINT ["java","-cp","app:app/lib/*","org/gatblau/onix/App"]
