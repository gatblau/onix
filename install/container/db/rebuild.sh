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

TAG=$1

if [ $# -eq 0 ]; then
    echo "An image tag is required for Onix. Provide it as a parameter."
    echo "Usage is: sh rebuild.sh [ONIX TAG]"
    exit 1
fi

# removes the container
docker rm -f onixdb

# deletes the image
docker rmi "onixdb:${TAG}"

# builds the image
sh build.sh "${TAG}"

# creates the container
sh run.sh "${TAG}"

# wait for the container to initialise
sleep 5

# shows the logs
docker logs onixdb
