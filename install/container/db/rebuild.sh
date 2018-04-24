#!/usr/bin/env bash
TAG=$1

if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh rebuild.sh [ONIX TAG]"
    exit 1
fi

# removes the container
docker rm -f onix-db

# deletes the image
docker rmi "onix-db:${TAG}"

# builds the image
sh build.sh "${TAG}"

# creates the container
sh run.sh "${TAG}"

# wait for the container to initialise
sleep 5

# shows the logs
docker logs onix-db
