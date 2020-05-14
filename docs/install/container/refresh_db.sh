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
command -v psql >/dev/null 2>&1 || { echo >&2 "psql is required but it's not installed. Aborting."; exit 1; }
command -v docker >/dev/null 2>&1 || { echo >&2 "docker is required but it's not installed. Aborting."; exit 1; }

# configuration
export PGPASSWORD=onix

HOST=localhost
PORT=5432
DB=onix
DBUSER=onix
DBPWD=onix
SPATH=${HOME}"/dev/onix/wapi/src/main/resources/db/install/4"

# removes the container
docker rm -f oxdb

# re-creates the container
docker run --name oxdb -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=${PGPASSWORD} "centos/postgresql-12-centos7"

# wait for the container to initialise
sleep 5

# shows the logs
docker logs oxdb

# re-deploys the database
psql -h ${HOST} -U postgres -c "CREATE DATABASE "${DB}";"
psql -h ${HOST} -U postgres -c "CREATE USER "${DBUSER}" WITH PASSWORD '"${DBPWD}"';"
psql -h ${HOST} -U postgres ${DB} -c "CREATE EXTENSION IF NOT EXISTS hstore;"
psql -h ${HOST} -U postgres ${DB} -c "CREATE EXTENSION IF NOT EXISTS intarray;"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/tables.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/json.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/validation.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/set.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/get.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/delete.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/queries.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/tree.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/tags.sql"
psql -h ${HOST} -U ${DBUSER} ${DB} -f "${SPATH}/keyman.sql"



