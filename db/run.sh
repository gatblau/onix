#!/usr/bin/env bash
# creates a docker container with the ONIX database
# usage:  sh run.sh admin_password
#
if [ ! -n "$1" ]; then
    admin_pwd="0n1x-8de5"
else
    admin_pwd=$1
fi

docker run --name onix-db -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=$admin_pwd onix-db:0.0.1-0
