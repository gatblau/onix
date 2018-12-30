#!/usr/bin/env bash
# gets the folder this script is in
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo '>>> Creating the ONIX user <<<'
psql -U postgres -c "CREATE USER onix WITH PASSWORD 'onix';"

echo '>>> Creating the ONIX database <<<'
createdb -E UTF8 -e -O onix onix

echo '>>> Installing HSTORE & UUID extensions <<<'
psql -U postgres -d onix -c 'create extension if not exists hstore;'
psql -U postgres -d onix -c 'create extension if not exists "uuid-ossp";'

echo '>>> Creating the database tables <<<'
psql -U postgres -d onix -a -f $DIR/sql/tables.sql

echo '>>> Creating the database functions <<<'
psql -U postgres -d onix -a -f $DIR/sql/validation_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/set_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/get_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/delete_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/queries.sql

echo '>>> ONIX database initialisation complete! <<<'