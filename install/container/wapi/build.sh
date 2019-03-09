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
    echo "Usage is: sh build.sh [ONIX TAG]"
    exit 1
fi

#../s2i build https://github.com/gatblau/onix/tree/v1 fabric8/s2i-java "gatoazul/onix-wapi:${TAG}"

# deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

echo "removes the target directory"
rm -rf ././../../../target/

echo "deletes the app temp folder"
rm -rf ./tmp

echo "packages the application"
mvn -f ././../../../pom.xml package

echo "unzips the application jar file"
unzip -o ././../../../target/*.jar -d ./tmp

echo "builds the docker image"
docker build -t creoworks/onixwapi-snapshot:${TAG} .

echo "tags the image as latest"
docker tag creoworks/onixwapi-snapshot:${TAG} creoworks/onixwapi-snapshot:latest