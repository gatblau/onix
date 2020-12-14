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
if [[ -z "${PUSH_IMAGE_REGISTRY+x}" ]]; then
    echo "PUSH_IMAGE_REGISTRY must be provided"
    exit 1
fi
if [[ -z "${PUSH_IMAGE_REPO+x}" ]]; then
    echo "PUSH_IMAGE_REPO must be provided"
    exit 1
fi
if [[ -z "${PUSH_IMAGE_NAME+x}" ]]; then
    echo "PUSH_IMAGE_NAME must be provided"
    exit 1
fi
if [[ -z "${PUSH_IMAGE_VERSION+x}" ]]; then
    echo "PUSH_IMAGE_VERSION must be provided"
    exit 1
fi
if [[ -z "${PULL_IMAGE_REGISTRY+x}" ]]; then
    echo "PULL_IMAGE_REGISTRY must be provided"
    exit 1
fi
if [[ -z "${BUILD_FILE_URL+x}" ]]; then
    echo "BUILD_FILE_URL must be provided"
    exit 1
fi
if [[ -z "${DOCKER_FILE_URL+x}" ]]; then
    echo "DOCKER_FILE_URL must be provided"
    exit 1
fi

# defines the container image fqn
PUSH_IMAGE_FQN="${PUSH_IMAGE_REGISTRY}/${PUSH_IMAGE_REPO}/${PUSH_IMAGE_NAME}"

# fetch the build.yaml file from the project repository
# if an authentication token has been provided
if [[ -z "${GIT_TOKEN+x}" ]]; then
  echo GIT_TOKEN not defined, retrieving build.yaml without authenticating
  wget "${BUILD_FILE_URL}" -O build.yaml
else
  echo GIT_TOKEN defined
  echo retrieving build.yaml with token
  wget --header="PRIVATE-TOKEN:${GIT_TOKEN}" "${BUILD_FILE_URL}" -O build.yaml
  echo retrieving Dockerfile with token
  wget --header="PRIVATE-TOKEN:${GIT_TOKEN}" "${DOCKER_FILE_URL}" -O Dockerfile
fi

# import required keys
# one or more public keys for verifying artefacts
# one private key for signing the container image
artie run import-keys

# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to import keys"
   exit 1
fi

# open artefact(s) using artie
artie run open-artefacts

# if the exit code of the previous command is not zero
if [ "$?" -ne 0 ]; then
   echo "failed to open artefact(s)"
   exit 1
fi

# login to the base image registry if a username is provided
if [[ -n "${PULL_IMAGE_REGISTRY_UNAME}" ]]; then
  # must have a password too
  if [[ -z "${PULL_IMAGE_REGISTRY_PWD+x}" ]]; then
    echo "PULL_IMAGE_REGISTRY_PWD must be provided"
    exit 1
  fi
  # login to the pull image registry (the one containing base image)
  buildah login -u "${PULL_IMAGE_REGISTRY_UNAME}" -p "${PULL_IMAGE_REGISTRY_PWD}" "${PULL_IMAGE_REGISTRY}"
fi

# performs the build of the new image defined by Dockerfile downloaded by init.sh
buildah --storage-driver vfs bud --isolation chroot -t "${PUSH_IMAGE_FQN}:${PUSH_IMAGE_VERSION}" .

# if there is a request to add latest tag
if [[ "${PUSH_IMAGE_TAG_LATEST}" = true ]]; then
  # tag the local image for docker.io
  buildah tag "${PUSH_IMAGE_FQN}:${PUSH_IMAGE_VERSION}" "${PUSH_IMAGE_FQN}:latest"
fi

# if credentials are provided to push the image to the push image registry
if [[ -n ${PUSH_IMAGE_REGISTRY_UNAME} ]]; then
  # must have a password too
  if [[ -z ${PUSH_IMAGE_REGISTRY_PWD+x} ]]; then
    echo "PUSH_IMAGE_REGISTRY_PWD must be provided"
    exit 1
  fi
  # push the new image to the registry using credentials
  buildah --storage-driver vfs push --tls-verify=false --creds "${PUSH_IMAGE_REGISTRY_UNAME}:${PUSH_IMAGE_REGISTRY_PWD}" "${PUSH_IMAGE_FQN}:${PUSH_IMAGE_VERSION}"

  # if required push the latest tag
  if [[ "${PUSH_IMAGE_TAG_LATEST}" = true ]]; then
     buildah --storage-driver vfs push --tls-verify=false --creds "${PUSH_IMAGE_REGISTRY_UNAME}:${PUSH_IMAGE_REGISTRY_PWD}" "${PUSH_IMAGE_FQN}:latest"
  fi
else
  # configure the push to use secrets in .dockercfg
  if [[ ! -f /var/run/secrets/openshift.io/push/.dockercfg ]]; then
    echo "ATTENTION! this build has been configured to use a docker-registry secret
the build process could not find .dockercfg in the container's file system, try the following:
  - ensure 'output/to/kind: DockerImage' is set within your build configuration so that the secret is added to the container file system
  - ensure you have configured a valid docker-registry secret
  or:
  - change the configuration to use push credentials instead
    do this by setting the PUSH_IMAGE_REGISTRY_UNAME and PUSH_IMAGE_REGISTRY_PWD variables"
    exit 1
  fi

  # buildah requires a slight modification to the push secret provided by the service
  # account in order to use it for pushing the image
  cp /var/run/secrets/openshift.io/push/.dockercfg /tmp
  (echo "{ \"auths\": " ; cat /var/run/secrets/openshift.io/push/.dockercfg ; echo "}") > /tmp/.dockercfg

  # push the new image to the registry using secrets
  buildah --storage-driver vfs push --tls-verify=false --authfile /tmp/.dockercfg "${PUSH_IMAGE_FQN}:${PUSH_IMAGE_VERSION}"
  # if required push the latest tag
  if [[ "${PUSH_IMAGE_TAG_LATEST}" = true ]]; then
    buildah --storage-driver vfs push --tls-verify=false --authfile /tmp/.dockercfg "${PUSH_IMAGE_FQN}:latest"
  fi
fi