# Building the docker images [(index)](./../readme.md)

In order to build the Onix docker images follow the steps below:

- Clone the [Onix](https://github.com/gatblau/onix.git) repository.
- Change the directory to the folder where [build.sh](../install/container/build.sh) script is.
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

**NOTE**: a new docker image tag is created automatically by the [build.sh](../install/container/build.sh) script and used for both the database and the service images.
The convention for the tag is as follows: **[abbreviated last commit hash].[ddMMyy-HHmmss]**