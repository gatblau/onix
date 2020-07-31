#!/usr/bin/env bash
#
#   Onix Service Status SeS - Copyright (c) 2018-2020 by www.gatblau.org
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
#   Unless required by applicable law or agreed to in writing, software distributed under
#   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
#   either express or implied.
#   See the License for the specific language governing permissions and limitations under the License.
#
#   Contributors to this project, hereby assign copyright in this code to the project,
#   to be licensed under the same terms as the rest of the code.
#
VERSION=$1
if [ $# -eq 0 ]; then
    echo "An image version is required for Onix DbMan. Provide it as a parameter."
    echo "Usage is: sh build.sh [APP VERSION] - e.g. sh build.sh v1.0.0"
    exit 1
fi

rm version

# creates a TAG for the newly built docker images
DATE=`date '+%d%m%y%H%M%S'`
HASH=`git rev-parse --short HEAD`
TAG="${VERSION}-${HASH}-${DATE}"

echo ${TAG} >> version

echo "TAG is: ${TAG}"

sleep 2