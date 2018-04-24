#!/usr/bin/env bash
TAG=$1
if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh build.sh [ONIX TAG]"
    exit 1
fi
../s2i build https://github.com/gatblau/onix.git fabric8/s2i-java "gatoazul/onix-svc:${TAG}"
