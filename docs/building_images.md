# Building the docker images

In order to build the Onix docker images follow the steps below:

- Clone the [Onix](https://github.com/gatblau/onix.git) repository.
- Change the directory to the folder where [build.sh](../install/container/build.sh) file is.
- Ensure [Docker](https://www.docker.com/) is installed on the host.
- Ensure the [s2i tool](https://github.com/openshift/source-to-image/releases) is installed on the folder.
- Execute the command below.

```bash
$ sh build.sh
```

Then check the images have been created:

```bash
$ docker images
```
