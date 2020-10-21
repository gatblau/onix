#!/usr/bin/env bash
#
#    Onix Config Manager - Copyright (c) 2018-2020 by www.gatblau.org
#
#    Licensed under the Apache License, Version 2.0 (the "License");
#    you may not use this file except in compliance with the License.
#    You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#    Unless required by applicable law or agreed to in writing, software distributed under
#    the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#    either express or implied.
#    See the License for the specific language governing permissions and limitations under the License.
#
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#
#  Usage: sh ox_up.sh
#  Launches Onix Web API and its backing PosgreSQL database in containers using Docker
#

echo check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

echo deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

echo try and delete existing Onix containers
docker rm -f oxdb
docker rm -f ox
docker rm -f dbman

echo "checking port 8080 is available for the Onix Web API process"
lsof -i:8080 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 8080 is in use, cannot continue: ensure there is no process using 8080/tcp port
  exit 1
fi

echo "checking port 5432 is available for the Onix database process"
lsof -i:5432 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 5432 is in use, cannot continue: ensure there is no process using 5432/tcp port
  exit 1
fi

echo "checking port 8085 is available for the Onix DbMan process"
lsof -i:8085 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 8085 is in use, cannot continue: ensure there is no process using 8085/tcp port
  exit 1
fi

echo creates the Onix database
docker run --name oxdb -it -d -p 5432:5432 \
    -e POSTGRESQL_ADMIN_PASSWORD=onix \
    "centos/postgresql-12-centos8"

echo creates the Onix DbMan
docker run --name dbman -itd -p 8085:8085 --link oxdb \
  -e OX_DBM_DB_HOST=oxdb \
  -e OX_DBM_DB_ADMINPWD=onix \
  -e OX_DBM_HTTP_AUTHMODE=none \
  -e OX_DBM_APPVERSION=0.0.4 \
  "gatblau/dbman-snapshot"

echo creates the Onix Web API
docker run --name ox -it -d -p 8080:8080 --link oxdb \
    -e DB_HOST=oxdb \
    -e DB_ADMIN_PWD=onix \
    -e WAPI_AUTH_MODE=basic \
    -e WAPI_ADMIN_USER=admin \
    -e WAPI_ADMIN_PWD=0n1x \
    "gatblau/onix-snapshot"

echo "please wait for the Web API to become available"
until $(curl --output /dev/null --silent --head --fail http://localhost:8080); do
    printf '.'
    sleep 2
done

echo "deploying database schemas"
echo "? creating the database"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1

echo "? deploying the schemas and functions"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1

echo 
echo "Web API ready to use @ localhost:8080"