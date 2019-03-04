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

# deletes any images with no tag
images_with_no_tag=$(docker images -f dangling=true -q)
if [ -n "$images_with_no_tag" ]; then
    docker rmi $images_with_no_tag
fi

# builds an onix postgresql database image using the S2I tool (https://github.com/openshift/source-to-image/releases)
../s2i build ./image_conf/ centos/postgresql-10-centos7 "creoworks/onixdb:${TAG}"

# https://github.com/sclorg/postgresql-container/tree/generated/10
