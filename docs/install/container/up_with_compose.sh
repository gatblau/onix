#!/usr/bin/env bash
#
#    Onix Config Manager - Copyright (c) 2018-2021 by www.gatblau.org
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
#

# loading environment vars
export $(grep -v '^#' .env)

echo "starting Onix with docker compose"
docker-compose up -d

echo "please wait..."
sleep 3

echo "? creating the database"
curl -H "Content-Type: application/json" \
     -H "Authorization: Basic $(printf '%s:%s' $DBMAN_HTTP_USER $DBMAN_HTTP_PWD | base64)" \
     -X POST http://localhost:8085/db/create 2>&1

echo "? deploying the schemas and functions"
curl -H "Content-Type: application/json" \
     -H "Authorization: Basic $(printf '%s:%s' $DBMAN_HTTP_USER $DBMAN_HTTP_PWD | base64)" \
     -X POST http://localhost:8085/db/deploy 2>&1

echo "? updating default's Onix Web API reader password"
until contents=$(curl -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" -H "Content-Type: application/json" -X PUT http://localhost:8080/user/$ONIX_HTTP_READER_USER/pwd -d "{\"pwd\":\"$ONIX_HTTP_READER_PWD\"}")
do
  sleep 3
done

echo "? updating default's Onix Web API writer password"
until contents=$(curl -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" -H "Content-Type: application/json" -X PUT http://localhost:8080/user/$ONIX_HTTP_WRITER_USER/pwd -d "{\"pwd\":\"$ONIX_HTTP_WRITER_PWD\"}")
do
  sleep 3
done

# NOTE: this has to be done at last as otherwise authentication will fail
echo "? updating default's Onix Web API admin password"
until contents=$(curl -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" -H "Content-Type: application/json" -X PUT http://localhost:8080/user/$ONIX_HTTP_ADMIN_USER/pwd -d "{\"pwd\":\"$ONIX_HTTP_ADMIN_PWD\"}")
do
  sleep 3
done