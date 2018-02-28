#!/usr/bin/env bash
# gets the folder this script is in
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo '>>> Creating the ONIX database <<<'
createdb onix

echo '>>> Creating the ONIX user <<<'
psql -U postgres -c "CREATE USER onix WITH PASSWORD 'onix';"

echo '>>> Creating the ONIX tables <<<'
psql -U postgres -d onix -a -f $DIR/create_tables.sql

echo '>>> ONIX database initialisation complete! <<<'