#!/usr/bin/env bash

VERSION=$1
if [ $# -eq 0 ]; then
    echo "An image version is required for Onix. Provide it as a parameter."
    echo "Usage is: sh build.sh [APP VERSION] - e.g. sh build.sh v1.0.0"
    exit 1
fi

# removes existing version file (if any)
rm src/main/resources/version

# creates a TAG for the newly built docker images
DATE=`date '+%d%m%y%H%M%S'`
HASH=`git rev-parse --short HEAD`
TAG="${VERSION}-${HASH}-${DATE}"

echo ${TAG} >> src/main/resources/version

echo "TAG is: ${TAG}"