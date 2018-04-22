#!/usr/bin/env bash

docker-compose down

# check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

# check docker-compose is installed
if [ ! -x "$(command -v docker-compose)" ]; then
    echo "Docker Compose is required to execute this script."
    exit 1
fi

# check s2i command is installed
if [ ! -x "$(command -v ./s2i)" ]; then
    echo "s2i is required to execute this script. See here: https://github.com/openshift/source-to-image"
    exit 1
fi

# check the onix db image exists in the registry
if [ $(docker inspect --type=image onix-db:0.0.1-0 2>&1 | grep -c "Error") -eq 1 ]; then
    cd db
    sh ./build.sh
    cd ..
fi

# check the onix service image exists in the registry
if [ $(docker inspect --type=image onix-svc:0.0.1-0 2>&1 | grep -c "Error") -eq 1 ]; then
    cd svc
    sh ./build.sh
    cd ..
fi

docker-compose up