#!/bin/bash
#
#    Sample Nexus script
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# This script is an example of how to spin up a Nexus OSS container
# along with any manual steps you need to perform in order to use it
# as a backend for Artisan Repository

. .env

NEXUS_PORT=8081
NEXUS_VOL=nexus_data

# Kill any old container and data
docker rm -f nexus
docker volume rm ${NEXUS_VOL}

# Start clean Nexus
docker run -d \
    -p ${NEXUS_PORT}:8081 \
    --name nexus \
    --restart always \
    --network ${DOCKER_NETWORK} \
    sonatype/nexus3

echo Waiting for Nexus to start up and configure itself ...
sleep 60

PASS=$(docker container exec nexus cat ./nexus-data/admin.password)

echo Next steps\:
echo - take a note of the admin password above
echo - browse to http://localhost:8081 and sign in as 'admin' with a password of ${PASS}
echo - go through the Nexus wizard as normal
echo - create a Repository of type 'raw (hosted)' named 'artisan'
