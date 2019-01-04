#!/usr/bin/env bash
#/*
#Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org
#
#Licensed under the Apache License, Version 2.0 (the "License");
#you may not use this file except in compliance with the License.
#You may obtain a copy of the License at
#
#http://www.apache.org/licenses/LICENSE-2.0
#
#Unless required by applicable law or agreed to in writing, software
#distributed under the License is distributed on an "AS IS" BASIS,
#WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#See the License for the specific language governing permissions and
#limitations under the License.
#
#Contributors to this project, hereby assign copyright in their code to the
#project, to be licensed under the same terms as the rest of the code.
#*/

# ----------------------------------------------------
# Name: prepare_onix_db.sh
# Description:
#   Creates user, database and tables required by the Onix CMDB application
#   Requires the pgsql client tool
#
# Usage:
#   $ sh configure_pgsql.sh localhost 5432 onix
# ----------------------------------------------------

# the database server hostname
DB_HOST=$1

# the database server port
DB_PORT=$2

# the password of the onix user
DB_PWD=$3

echo 'creating the database user...'
psql -U postgres -h $DB_HOST -p $DB_PORT -c "CREATE USER onix WITH PASSWORD '$DB_PWD';"

echo 'creating the onix database...'
createdb -h $DB_HOST -p $DB_PORT -E UTF8 -e -O onix onix

echo '>>> Installing HSTORE extension <<<'
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -c 'create extension if not exists hstore;'

echo '>>> Creating the database objects <<<'
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/create_tables.sql
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/validation_funcs.sql
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/set_funcs.sql
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/get_funcs.sql
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/delete_funcs.sql
psql -U postgres -h $DB_HOST -p $DB_PORT -d onix -a -f ../container/db/image_conf/postgresql-init/queries.sql

