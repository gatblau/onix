#!/usr/bin/env bash
TAG=$1
../s2i build https://github.com/gatblau/onix.git fabric8/s2i-java "gatoazul/onix-svc:${TAG}"
