#!/usr/bin/env bash

# check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

# check s2i command is installed
if [ ! -x "$(command -v ./s2i)" ]; then
    echo "s2i is required to execute this script. See here: https://github.com/openshift/source-to-image"
    exit 1
fi

# creates a TAG for the newly built docker images
DATE=`date '+%d%m%y-%H%M%S'`
HASH=`git rev-parse --short HEAD`
ONIXTAG="${HASH}.${DATE}"
echo "Onix TAG is: ${ONIXTAG}"

# builds the onix-db image
cd db
sh ./build.sh $ONIXTAG
cd ..

# builds the onix-wapi image
cd wapi
sh ./build.sh $ONIXTAG
cd ..