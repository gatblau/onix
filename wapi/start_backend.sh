#!/usr/bin/env bash
#
#    Onix Config Manager - Web API Backend services for testing
#    Copyright (c) 2018-2020 by www.gatblau.org
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
#  Usage: sh start_backend.sh
#  Description: launches 2 containers using Docker in the local manchine:
#   - PostgreSQL (TCP 5432) backend database for Onix Web API
#   - Artemis (TCP 8161) backend AMQP broker for Onix Web API
#

echo check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

echo deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

echo try and delete existing backend containers
docker rm -f oxmsg
docker rm -f oxdb
docker rm -f ox

echo "checking port 8080 is available for the Onix Web API"
lsof -i:8080 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 8080 is in use, cannot continue: ensure there is no process using 8080/tcp port
  exit 1
fi

echo "checking port 8161 is available for the Artemis broker process"
lsof -i:8161 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 8161 is in use, cannot continue: ensure there is no process using 8161/tcp port
  exit 1
fi

echo "checking port 61616 is available for the Artemis broker process"
lsof -i:61616 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 61616 is in use, cannot continue: ensure there is no process using 61616/tcp port
  exit 1
fi

echo "checking port 5432 is available for the Onix database process"
lsof -i:5432 | grep LISTEN
RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo port 5432 is in use, cannot continue: ensure there is no process using 5432/tcp port
  exit 1
fi

echo creates the Onix database container
docker run --name oxdb -it -d -p 5432:5432 \
    -e POSTGRESQL_ADMIN_PASSWORD=onix \
    "centos/postgresql-12-centos7"

echo creates the Onix message broker container
docker run --name oxmsg -it -d \
  -e ARTEMIS_USERNAME=admin \
  -e ARTEMIS_PASSWORD=adm1n \
  -p 8161:8161 \
  -p 61616:61616 \
  vromero/activemq-artemis
