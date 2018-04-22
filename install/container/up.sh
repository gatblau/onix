#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh up.sh [ONIX TAG]"
    exit 1
fi

# assigns the first parameter to the ONIXTAG
ONIXTAG=$1

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

# exports ONIXTAG as an environment variable so that it is available within the docker-compose.yml
export ONIXTAG

docker-compose up #--detach
