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
# re-creates a docker container with a postgres database
# usage:  sh refresh.sh
#
# removes the container
docker rm -f onixdb

# re-creates the container
docker run --name onixdb -it -d -p 5432:5432 -e POSTGRESQL_ADMIN_PASSWORD=onix "centos/postgresql-10-centos7"

# wait for the container to initialise
sleep 5

# shows the logs
docker logs onixdb

# re-deploys the database
curl localhost:8080/ready
