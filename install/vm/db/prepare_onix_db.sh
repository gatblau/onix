#!/usr/bin/env bash
#
# Name: prepare_onix_db.sh
# Description:
#   Creates user, database and tables required by the Onix CMDB application
#   Requires the pgsql client tool
#
# Usage:
#   $ sh configure_pgsql.sh localhost 5432 onix
# ----------------------------------------------------
#
# the database server hostname
HOST=$1

# the database server port
PORT=$2

# the password of the onix user
DB_PWD=$3

echo 'creating the database user...'
psql -U postgres -h $HOST -p $PORT -c "CREATE USER onix WITH PASSWORD '$DB_PWD';"

echo 'creating the onix database...'
createdb -h $HOST -p $PORT -E UTF8 -e -O onix onix

echo 'creating the db tables...'
psql -U postgres -h $HOST -p $PORT -d onix -a -f ../container/db/image_conf/postgresql-init/create_tables.sql