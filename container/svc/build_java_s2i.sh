#!/usr/bin/env bash

JAVA_DOCKER="ocp_s2i_java"
JAVA_JAR="onix-1.0-SNAPSHOT.jar"
APP_REPO="onix"

if [ ! -d "$JAVA_DOCKER" ]; then
    echo Retrieving Java S2I builder image source code
    git clone https://github.com/gatblau/$JAVA_DOCKER.git
    cd $JAVA_DOCKER/

    echo Building Java S2I builder image
    docker build -t java:1.0 .
    cd ..
fi