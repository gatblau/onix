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
# re-creates a docker container with a postgres database for testing only
# usage:  sh refresh.sh

# PSQL install notes - RHEL
# sudo yum install https://download.postgresql.org/pub/repos/yum/10/redhat/rhel-7-x86_64/pgdg-redhat10-10-2.noarch.rpm
# sudo yum install postgresql10

# PSQL install notes - MacOS
# brew install libpq
# brew link --force libpq

# pre-requisites
command -v docker >/dev/null 2>&1 || { echo >&2 "docker is required but it's not installed. Aborting."; exit 1; }

APP_VER="0.0.4"
DBPWD=onix

docker rm -f oxdb
docker rm -f dbman

echo "? starting a new database container"
docker run --name oxdb -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=${DBPWD} "centos/postgresql-12-centos8"

echo "? waiting for the database to start before proceeding"
sleep 5

echo "? launching DbMan"
docker run --name dbman -itd -p 8085:8085 --link oxdb \
  -e OX_DBM_DB_HOST=oxdb \
  -e OX_DBM_DB_ADMINPWD=${DBPWD} \
  -e OX_DBM_HTTP_AUTHMODE=none \
  -e OX_DBM_APPVERSION=${APP_VER} \
  "gatblau/dbman-snapshot"

echo "? please wait for DbMan to become available"
sleep 3

echo "? creating the database"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1

echo "? deploying the schemas and functions"
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1
