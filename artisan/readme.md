<img src="https://github.com/gatblau/artisan/raw/master/artisan.png" width="150" align="right"/>

# ARTISAN - One stop shop for building, packaging, signing, tagging and distributing any application.

Artisan is part of Onix Configuration Manager build system, an effort to standardise the process of building and packaging software to run in container platforms, so that tasks such as automated patching of security vulnerabilities can be realised.

Artisan is a command line tool like docker that allows you to:

- build any type of application (e.g. java, node/javascript, golang, etc)
- package any application in a consistent way, and prepare it for easily embedding into container images
- automatically add metadata and digital signatures to the packages
- tag artefacts as you do with container images
- push and pull packages to and from artefact registries
- open packages automatically verifying their digital signature and prompting for tampering
- use it in linux, windows, and mac os
- run it as an artefact registry which connects to backends such as Nexus, Artifactory, S3 and file systems.
- facilitate creation of automation commands by allowing creation of complex command-based functions and calling functions within functions

## Why decoupling the building of applications from the creation of container images?

1. So that automated runtime patching of production container images can be facilitated
2. So that applications do not have to go through complex build processes every time their base image has to be updated
3. So that rebuilding container images is faster, particularly when multiple artefacts are required by the image build
4. More efficient utilisation of a container platform CPU cycles
5. Increased speed of delivering security updates to production by reducing the risk introduced by changes
6. So that the application build pipeline does not have to be invoked if a patch has to be applied to the application base image and the application has not changed
7. So that there is an easy way to distribute any command tools built on any language and deployed on any platform (e.g. CLIs)

## Building, packaging, and distributing artefacts

Building an artefact is as easy as running the command below, Artisan relies on a [build.yaml](build.yaml) file in the git project root.

```sh
# build from a git repo
# NOTE: you will need the golang toolchain installed on your machine as boot is built using go
./art build -t localhost:8081/gatblau/boot:v01 https://github.com/gatblau/boot

# see the created artefact
./art list

# see artefact numeric ids only
./art list -q

# push to a remote artefact registry with a Nexus 3 backend
./art push localhost:8082/gatblau/boot:v01 -u="admin:admin123" -t=false

# delete the artefact from the local artefact registry using the name:tag
./art rm localhost:8081/gatblau/boot:v01

# see the updated local artefacts list
./art list

# pull the artefact from the remote artefact registry
./art pull localhost:8081/gatblau/boot:v01

# see the updated local artefacts list
./art list

# delete all artefacts
./art rm $(./art list -q)
```

## Building Artisan with Artisan

If you want to build artisan with artisan do the following:

```sh
# NOTE: you will need the golang toolchain installed on your machine
./art build -t gatblau/artisan -f artisan https://github.com/gatblau/onix
```

## Implemented backends

At the moment the only backend implemented is [Sonatype Nexus 3](https://help.sonatype.com/repomanager3)
Other backends are likely to be added in the future.

Artisan can be launched as an artefact registry as follows:

```sh
./art serve
```

For a comprehensive list of the available commands:

```sh
./art --help
```