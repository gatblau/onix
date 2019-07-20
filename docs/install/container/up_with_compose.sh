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
echo "starting Onix with docker compose"
docker-compose up -d

echo "please wait for the Web API to become available"
sleep 10

echo "deploying database schemas"
curl localhost:8080/ready

echo 
echo "Web API ready to use @ localhost:8080"
echo "Ox Kube ready to use @ localhost:8000"
