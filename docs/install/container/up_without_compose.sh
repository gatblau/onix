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
#  Launches Onix in the local machine without using docker-compose
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
docker rm -f dbman
docker rm -f ox
docker rm -f oxku

echo creates the Onix database
docker run --name oxdb -it -d -p 5432:5432 \
    -e POSTGRESQL_ADMIN_PASSWORD=onix \
    "centos/postgresql-12-centos7"

echo creates DbMan
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

sleep 5

# uncomment below to deploy additional services

#echo create the Onix Kubernetes agent
#docker run --name oxku -it -d -p 8000:8000 --link ox \
#    -e OXKU_ID=kube-01 \
#    -e OXKU_ONIX_URL=http://ox:8080 \
#    -e OXKU_CONSUMERS_CONSUMER=webhook \
#    -e OXKU_LOGINLEVEL=Trace \
#    -e OXKU_ONIX_AUTHMODE=basic \
#    -e OXKU_ONIX_USER=basic \
#    -e OXKU_ONIX_PASSWORD=0n1x \
#    "gatblau/oxkube-snapshot"
#
#echo create the Onix Web Console
#docker run --name oxwc -it -d -p 3000:3000 --link ox \
#    -e WC_OX_WAPI_URI=http://onix:8080 \
#    -e WC_OX_WAPI_AUTH_MODE=basic \
#    "gatblau/oxwc-snapshot"
#
#echo "please wait for the Web API to become available"
#sleep 10

echo "? creating the database"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1

echo "? deploying the schemas and functions"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1