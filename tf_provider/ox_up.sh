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
#  Usage: sh ox_up.sh
#  Launches Onix Web API and its backing PosgreSQL database in containers using Docker
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

echo try and delete existing Onix containers
docker rm -f oxdb
docker rm -f ox

echo creates the Onix database
docker run --name oxdb -it -d -p 5432:5432 \
    -e POSTGRESQL_ADMIN_PASSWORD=onix \
    "centos/postgresql-12-centos7"

echo creates the Onix Web API
docker run --name ox -it -d -p 8080:8080 --link oxdb \
    -e DB_HOST=oxdb \
    -e DB_ADMIN_PWD=onix \
    -e WAPI_AUTH_MODE=basic \
    -e WAPI_ADMIN_USER=admin \
    -e WAPI_ADMIN_PWD=0n1x \
    "gatblau/onix-snapshot"

echo "please wait for the Web API to become available"
until $(curl --output /dev/null --silent --head --fail http://localhost:8080); do
    printf '.'
    sleep 2
done

echo "deploying database schemas"
curl localhost:8080/ready

echo 
echo "Web API ready to use @ localhost:8080"