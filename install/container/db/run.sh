#!/usr/bin/env bash
# creates a docker container with the ONIX database
# usage:  sh run.sh tag
#
TAG=$1

if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh run.sh [ONIX TAG]"
    exit 1
fi

docker run --name onix-db -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=onix "gatoazul/onix-db:${TAG}"
