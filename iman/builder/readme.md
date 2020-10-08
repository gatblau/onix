# Image Manager Builder

The image the Image Manager uses to build container images from:

  1. Application package
  2. Base Image
  3. Build Artefacts

## Building the builder image

```sh
$ make set-version
$ make snapshot-image
```

## How it works

The builder image uses [buildah](https://buildah.io/) to create container images.

When the image is run, the bootstrap script [build.sh](./build.sh) is executed.

This script downloads a initialisation script from a remote git repository and executes it.

The initialisation script sole purpose is to fetch all the artefacts required by the builder to build the image.  

For example:

- the Dockerfile
- the application package (binaries)
- Any other required artefacts

Then it calls buildah to build the image using the downloaded dockerfile and required artefacts. 

Finally, buildah pushes the image to the image registry.

## Builder configuration

The builder image requires the following variables:

| var | description | required |
|---|---|---|
| **IMAGE_REGISTRY** | the image reqistry for the new image (e.g. quay.io) | yes |
| **IMAGE_REPO** | the image repository | yes |
| **IMAGE_NAME** | the image name | yes |
| **IMAGE_VERSION** | the image version tag | yes |
| **INIT_SCRIPT_URL** | the url for the init.sh script | yes |
| **IS_SNAPSHOT** | if defined, a "snapshot" tag is added to the image name. | no |
