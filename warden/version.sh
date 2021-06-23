#!/usr/bin/env bash
TAG=$1
if [ $# -eq 0 ]; then
    echo "A tag is required for Onix Warden. Provide it as a parameter."
    echo "Usage is: sh version.sh [TAG] - e.g. sh version.sh my-tag"
    exit 1
fi

rm ./core/version.go || true
rm ./version || true
pwd
printf "package core\nconst Version=\"%s\"" "${TAG}" > ./core/version.go
echo "${TAG}" >> version