<img src="./pics/ox.png" width="200" height="200" align="right">

# Building the docker image [(index)](./../readme.md)

The build process uses the [make](https://www.gnu.org/software/make/manual/make.html) utility. [Git]() and [Docker](https://docs.docker.com/engine/reference/commandline/cli/) are also required.

## Web API

To build the web api, no dependencies are required other than :

```bash
# clone the repository
git clone https://github.com/gatblau/onix

# navigate to the wapi root folder
cd ./onix/wapi

# create a version tag from commit and time stamp
make version

# build the docker image
make image

# check the onix image has been created
 docker images
```

## Ox-Kube

Building ox-kube is done in the same way as with the web api as follows:

```bash
# clone the repository
git clone https://github.com/gatblau/onix

# navigate to the oxkube root folder
cd ./onix/agents/oxkube

# create a version tag from commit and time stamp
make version

# build the docker image
make image

# check the oxkube image has been created
 docker images
```

## Version Tag

The "*make version*" command creates a tag following the convention:

**[semantic version number]-[git abbreviated last commit hash]-[ddMMyyHHmmss]**

For example: v0.0.2-b2a6da0-170719155635