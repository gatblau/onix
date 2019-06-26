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

# Builds the Onix image

# check docker is installed
if [ ! -x "$(command -v docker)" ]; then
    echo "Docker is required to execute this script."
    exit 1
fi

VERSION=$1
if [ $# -eq 0 ]; then
    echo "An image version is required for Onix. Provide it as a parameter."
    echo "Usage is: sh build.sh [ONIX VERSION] - e.g. sh build.sh v1.0.0"
    exit 1
fi

# creates a TAG for the newly built docker images
DATE=`date '+%d%m%y%H%M%S'`
HASH=`git rev-parse --short HEAD`
ONIXTAG="${VERSION}-${HASH}-${DATE}"

#ONIXTAG="${VERSION}"

echo "Onix TAG is: ${ONIXTAG}"

# writes the version to the resources folder so it is available for the wapi to report on
rm ../../src/main/resources/version
echo ${ONIXTAG} >> ../../src/main/resources/version

# removes existing images from the local registry
docker images -a | grep "gatblau*" | awk '{print $3}' | xargs docker rmi -f

# builds the onix image
echo "building onix image"

# deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

echo "removing the target directory"
rm -rf ././../../target/

echo "deleting the app temp folder"
rm -rf ./tmp

echo "packaging the application"
mvn -f ././../../pom.xml package

echo "unzipping the application jar file"
unzip -o ././../../target/*.jar -d ./tmp

echo "building the docker image"
#docker build -t gatblau/onix:$ONIXTAG .
docker build -t gatblau/onix-snapshot:$ONIXTAG .

echo "tagging the image as latest"
#docker tag gatblau/onix:$ONIXTAG gatblau/onix:latest
docker tag gatblau/onix-snapshot:$ONIXTAG gatblau/onix-snapshot:latest

echo "pushing to docker.io"
#docker push gatblau/onix:$ONIXTAG
docker push gatblau/onix-snapshot:$ONIXTAG
