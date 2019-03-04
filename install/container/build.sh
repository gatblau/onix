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

# Builds the Onix WAPI and DB images required to run the application

# check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

# check s2i command is installed
if [ ! -x "$(command -v ./s2i)" ]; then
    echo "s2i is required to execute this script. See here: https://github.com/openshift/source-to-image"
    exit 1
fi

# creates a TAG for the newly built docker images
DATE=`date '+%d%m%y-%H%M%S'`
HASH=`git rev-parse --short HEAD`
ONIXTAG="${HASH}-${DATE}"
echo "Onix TAG is: ${ONIXTAG}"

# builds the onix-db image
echo building onixdb...
cd db
sh ./build.sh $ONIXTAG
cd ..

# builds the onixwapi image
echo building onixwapi
cd wapi
sh ./build.sh $ONIXTAG
cd ..