#!/usr/bin/env bash
TAG=$1
# builds an onix postgresql database image using the S2I tool (https://github.com/openshift/source-to-image/releases)
../s2i build ./image_conf/ centos/postgresql-10-centos7 "gatoazul/onix-db:${TAG}"
