#!/usr/bin/env bash
# pass the following variables in...
HOST=$1
PORT=$2
DB_PWD=$3

echo 'creating the database user...'
psql -U postgres -h $HOST -p $PORT -c "CREATE USER onix WITH PASSWORD '$DB_PWD';"

echo 'creating the onix database...'
createdb -h $HOST -p $PORT -E UTF8 -e -O onix onix

echo 'creating the db tables...'
psql -U postgres -h $HOST -p $PORT -d onix -a -f ../container/db/image_conf/postgresql-init/create_tables.sql