#!/usr/bin/env bash
#
# Sentinel - Copyright (c) 2019 by www.gatblau.org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
# Unless required by applicable law or agreed to in writing, software distributed under
# the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
# either express or implied.
# See the License for the specific language governing permissions and limitations under the License.
#
# Contributors to this project, hereby assign copyright in this code to the project,
# to be licensed under the same terms as the rest of the code.
#
# Builds the Sentinel OCI image using buildah.
# Replaces the Dockerfile build
# Run in CentOS 8 / RHEL 8
#
# NOTE:
#  if running on a buildah version lower than 1.11.3, the run with sudo e.g. sudo bash build.sh $1 $2 $3 $4 $5 $6
#   https://cbs.centos.org/koji/rpminfo?rpmID=171623
#
# check the input variables
if [ $# -lt 6 ]; then
    echo "Parameters provided are not correct. See usage below."
    echo "Usage is: sh build.sh [REPO_NAME] [APP_NAME] [APP_VERSION] [SNAPSHOT] - e.g. sh build.sh gatblau sentinel v0.0.1 yes"
    exit 1
fi

# read the input variables
REPO_NAME=$1
APP_NAME=$2
APP_VERSION=$3
SNAPSHOT=$4
IMG_FMT=$5 # image format i.e. docker or oci
IMG_REG=$6 # image registry i.e. docker.io or quay.io

# if building snapshot images, append snapshot to the end of the APP_NAME
if [ $SNAPSHOT == "yes" ]; then
    echo "building snapshot images"
    APP_NAME=$APP_NAME"-snapshot"
else
    echo "building release images"
fi

# set the path to the go command within the golang container
go=/usr/local/go/bin/go

# create a golang builder working container to build Sentinel
builder=$(buildah from docker://golang)

# set the working directory outside of the $GOPATH so the build does not fail
buildah config --workingdir /app $builder

# mount the container file system
builder_mnt=$(buildah mount $builder)

# copy the source code to the container
buildah copy $builder . .

# get the dependencies
buildah run $builder $go get .

# configure the environment for the build
buildah config --env CGO_ENABLED=0 $builder
buildah config --env GOOS=linux $builder

# build the Sentinel binary
buildah run $builder $go build -o sentinel .

# create the Sentinel working container
sentinel=$(buildah from docker://registry.access.redhat.com/ubi8/ubi-minimal)

# set labels
buildah config --label maintainer="GATBLAU <sentinel@gatblau.org>" $sentinel
buildah config --label author="gatblau.org>" $sentinel

# set the working directory for the Sentinel binary
buildah config --workingdir /app $sentinel

# copy the Sentinel binary and config file from the go builder mount to the Sentinel working container
buildah copy $sentinel $builder_mnt/app/sentinel ./
buildah copy $sentinel $builder_mnt/app/config.toml ./

# set the default user running containers based on this image
buildah config --user 20 $sentinel

# set the command to run Sentinel
buildah config --cmd "./sentinel" $sentinel

# commit the Sentinel working container to an image in the local registry
buildah commit --format $IMG_FMT $sentinel $REPO_NAME/$APP_NAME:$APP_VERSION

# tag the local image for the registry specified (i.e. docker.io or quay.io)
buildah tag $REPO_NAME/$APP_NAME:$APP_VERSION $IMG_REG/$REPO_NAME/$APP_NAME:$APP_VERSION
buildah tag $REPO_NAME/$APP_NAME:$APP_VERSION $IMG_REG/$REPO_NAME/$APP_NAME:latest

# remove the working containers
buildah rm $builder
buildah rm $sentinel
