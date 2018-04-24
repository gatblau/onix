#!/usr/bin/env bash
TAG=$1

if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh build.sh [ONIX TAG]"
    exit 1
fi

# builds an onix postgresql database image using the S2I tool (https://github.com/openshift/source-to-image/releases)
../s2i build ./image_conf/ centos/postgresql-10-centos7 "gatoazul/onix-db:${TAG}"
