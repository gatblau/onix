#!/usr/bin/env bash
TAG=$1
if [ $# -eq 0 ]; then
    echo "A tag is required for Onix Artisan. Provide it as a parameter."
    echo "Usage is: sh version.sh [TAG] - e.g. sh version.sh my-tag"
    exit 1
fi

rm ./docs/version.go || true
rm version || true
printf "package docs\nconst Version=\"%s\"" "${TAG}" >> ./docs/version.go
echo "${TAG}" >> version