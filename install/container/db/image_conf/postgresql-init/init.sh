#!/usr/bin/env bash
#
#    Onix CMDB - Copyright (c) 2018-2019 by www.gatblau.org
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

# gets the folder this script is in
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo '>>> Creating the ONIX user <<<'
psql -U postgres -c "CREATE USER onix WITH PASSWORD 'onix';"

echo '>>> Creating the ONIX database <<<'
createdb -E UTF8 -e -O onix onix

echo '>>> Installing HSTORE extension <<<'
psql -U postgres -d onix -c 'create extension if not exists hstore;'

echo '>>> Creating the database tables <<<'
psql -U postgres -d onix -a -f $DIR/sql/tables.sql

echo '>>> Creating the database functions <<<'
psql -U postgres -d onix -a -f $DIR/sql/validation_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/set_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/get_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/delete_funcs.sql
psql -U postgres -d onix -a -f $DIR/sql/queries.sql

echo '>>> ONIX database initialisation complete! <<<'