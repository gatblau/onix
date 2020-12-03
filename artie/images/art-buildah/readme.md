# Artie Linux Container Image Builder

A container image with **artie** and **buildah** that can be used to create container images out of:

- artefacts in an artefact registry
- base images in a container registry

## How it works

The builder image uses [buildah](https://buildah.io/) to create container images for applications.
The builder also has artie to allow pulling artefacts from the Artie artefact registry to embed in the container image.
Because the application is packaged by Artie this image does not need any tool chain to build the application,
that is already built and in the Artie artefact registry.

When the image is run, the bootstrap script [build.sh](./build.sh) is executed.

This script downloads a initialisation script from a remote git repository and executes it.

The initialisation script sole purpose is to fetch all the artefacts required by the builder to build the image.  

Typically:

- the Dockerfile
- the application artefacts - pulled using Artie CLI
- Any other required assets

Then it calls buildah to build the image using the downloaded dockerfile and required artefacts.

Finally, buildah pushes the image to the image registry.

## Builder configuration

The builder image requires the following variables:

| var | description | required |
|---|---|---|
| **PUSH_IMAGE_REGISTRY** | the image reqistry for the new image (e.g. quay.io) | yes |
| **PUSH_IMAGE_REPO** | the new image repository | yes |
| **PUSH_IMAGE_NAME** | the new image name | yes |
| **PUSH_IMAGE_VERSION** | the new image version tag | yes |
| **PUSH_IMAGE_REGISTRY_UNAME** | the username used to login to the new image registry.  | no. If not provided, credentials must be set in a docker-registry secret. |
| **PUSH_IMAGE_REGISTRY_PWD** | the password used to login to the new image registry | yes |
| **PUSH_IMAGE_TAG_LATEST** | whether to tag the new image as latest | no |
| **INIT_SCRIPT_URL** | the url for the init.sh script | yes |
| **GIT_TOKEN** | an authentication token for the git repository where the build scripts are. | no. If not provided, then no authentication is used. |
| **PULL_IMAGE_REGISTRY** | the registry where thw base image used for the build is | yes |
| **PULL_IMAGE_REGISTRY_UNAME** | the username used to login to the base image registry | no. If not provided, then it is assumed that the base image is public |
| **PULL_IMAGE_REGISTRY_PWD** | the password used to login to the base image registry | only if PULL_IMAGE_REGISTRY_UNAME is provided |