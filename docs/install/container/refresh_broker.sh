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
# launches Artemis in a container
#
# usage:  sh refresh_broker.sh

# removes the container
docker rm -f oxmsg

docker run --name oxmsg -it -d \
  -e ARTEMIS_USERNAME=amqpadmin \
  -e ARTEMIS_PASSWORD=amqppassw0rd \
  -p 8161:8161 \
  -p 61616:61616 \
  -p 5672:5672 \
  vromero/activemq-artemis