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
#  Launches Onix in the local machine without using docker-compose
#
# check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

# deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

# try and delete existing Onix containers
docker rm -f oxdb
docker rm -f ox

# creates the Onix containers
docker run --name oxdb -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=onix "centos/postgresql-10-centos7"
docker run --name ox -it -d -p 8080:8080 --link oxdb -eDB_HOST=oxdb "gatblau/onix-snapshot"