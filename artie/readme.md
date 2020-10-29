# ARTIE - Docker like artefact builder, packager and publisher.

Build and package application artefacts as you do with container images using a binary artefact repository.

Artie is designed to apply the docker experience to managing application artefacts.

## Building from a remote repository

```sh
# build from a git repo
./artie build -t localhost:8081/gatblau/boot:v01 https://github.com/gatblau/boot

# see created artefact
./artie artefacts

# see artefact numeric ids only
./artie artefacts -q

# push to Nexus 3 (gatblau raw repository) running on localhost:8081
./artie push localhost:8081/gatblau/boot:v01 -u "admin:admin123"

# delete using the name:tag
./artie rm localhost:8081/gatblau/boot:v01

# delete using artefact some characters in artifact Id
./artie rm 56734

# delete all artefacts
./artie rm $(./artie artefacts -q)
```

## Building Artie with Artie

If you want to build artie with artie do the following:

```sh
./artie build -t localhost:8081/gatblau/artie -f artie https://github.com/gatblau/onix
```

## Implemented backends

At the moment the only backend implemented is [Sonatype Nexus 3](https://help.sonatype.com/repomanager3)
Other backends are likely to be added in the future.