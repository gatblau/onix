#!/usr/bin/env bash
# gets the folder this script is in
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo '>>> Creating the ONIX user <<<'
psql -U postgres -c "CREATE USER onix WITH PASSWORD 'onix';"

echo '>>> Creating the ONIX database <<<'
createdb -E UTF8 -e -O onix onix

echo '>>> Installing HSTORE extension <<<'
psql -U postgres -d onix -c 'create extension hstore;'

echo '>>> Creating the ONIX tables <<<'
psql -U postgres -d onix -a -f $DIR/create_tables.sql

echo '>>> Creating the ONIX functions <<<'
psql -U postgres -d onix -a -f $DIR/create_funcs.sql

echo '>>> ONIX database initialisation complete! <<<'