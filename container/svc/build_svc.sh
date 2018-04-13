#!/usr/bin/env bash

JAVA_DOCKER="ocp_s2i_java"
JAVA_JAR="onix-1.0-SNAPSHOT.jar"
APP_REPO="onix"

if [ ! -f "$JAVA_JAR" ]; then
    echo building application
    git clone https://github.com/gatblau/$APP_REPO.git
    cd $APP_REPO/
    mvn package
    cp ./target/$JAVA_JAR ../$JAVA_JAR
    cd ..
    rm -rf $APP_REPO
fi

docker build -t onix-svc:1.0 .

