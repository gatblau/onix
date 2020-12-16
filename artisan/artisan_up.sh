#!/bin/bash
docker-compose up -d
echo waiting until Nexus is up
until $(curl --output /dev/null --silent --head --fail http://localhost:8081); do
    printf '.'
    sleep 5
done
# retrieves the current Nexus admin password
pwd=$(docker container exec nexus cat ./nexus-data/admin.password)
echo "the current credentials are : admin:${pwd}"
echo "1. log into nexus and update it to admin:admin"
echo "2. create a new raw (hosted) repository called 'artisan'"
