# ARTIE - Artefact Builder and Packager

Build and package application artefacts as you do with container images using a binary artefact repository.

## Building from a remote repository

```sh
# build from a git repo
./artie build -t localhost:8081/gatblau/boot:v01 https://github.com/gatblau/boot

# see created artefact
./artie artefacts

# push to Nexus 3 (gatblau raw repository) running on localhost:8081
./artie push localhost:8081/gatblau/boot:v01 -u "admin:admin123"
```