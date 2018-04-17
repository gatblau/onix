#!/usr/bin/env bash
# removes the container
docker rm -f onix-db

# deletes the image
docker rmi onix-db:0.0.1-0

# builds the image
sh build.sh

# creates the container
sh run.sh onix

# wait for the container to initialise
sleep 5

# shows the logs
docker logs onix-db
