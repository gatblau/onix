#!/bin/sh
#
#    Onix Config Manager - Image Manager - Image Builder
#    Copyright (c) 2018-2020 by www.gatblau.org
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

# custom builds: https://docs.openshift.com/container-platform/4.5/builds/custom-builds-buildah.html
# securing builds: https://docs.openshift.com/container-platform/4.5/builds/securing-builds-by-strategy.html#securing-builds-by-strategy

# pre-condition health check: all required environment variables have been provided!
if [[ -z ${IMAGE_REGISTRY+x} ]]; then
    echo "IMAGE_REGISTRY must be provided"
    exit 1
fi
if [[ -z ${IMAGE_REPO+x} ]]; then
    echo "IMAGE_REPO must be provided"
    exit 1
fi
if [[ -z ${IMAGE_NAME+x} ]]; then
    echo "IMAGE_NAME must be provided"
    exit 1
fi
if [[ -z ${IMAGE_VERSION+x} ]]; then
    echo "IMAGE_VERSION must be provided"
    exit 1
fi
if [[ -z ${INIT_SCRIPT_URL+x} ]]; then
    echo "INIT_SCRIPT_URL must be provided"
    exit 1
fi

# defines the container image fqn
IMAGE_FQN="${IMAGE_REGISTRY}/${IMAGE_REPO}/${IMAGE_NAME}"

# if the variable IS_SNAPSHOT is defined then add "snapshot" to the name
if [[ ! -z ${IS_SNAPSHOT+x} ]]; then
    IMAGE_FQN="${IMAGE_FQN}-snapshot"
fi

# fetch the init script
curl -o init.sh "${INIT_SCRIPT_URL}"

# run the script
sh init.sh

# performs the build of the new image defined by Dockerfile downloaded by init.sh
buildah --storage-driver vfs bud --isolation chroot -t "${IMAGE_FQN}:${IMAGE_VERSION}" .

# buildah requires a slight modification to the push secret provided by the service
# account in order to use it for pushing the image
cp /var/run/secrets/openshift.io/push/.dockercfg /tmp
(echo "{ \"auths\": " ; cat /var/run/secrets/openshift.io/push/.dockercfg ; echo "}") > /tmp/.dockercfg

# push the new image to the target for the build
buildah --storage-driver vfs push --tls-verify=false --authfile /tmp/.dockercfg "${IMAGE_FQN}:${IMAGE_VERSION}"