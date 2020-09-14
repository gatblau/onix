#!/bin/bash

echo ==================================================
echo Starting Onix pod ...
echo ==================================================
podman pod create --hostname oxpod --name oxpod \
  -p 80:80,8080:8080,8085:8085,1883:1883

echo ==================================================
echo Starting Onix database ...
echo ==================================================
podman run -d --name oxdb \
  --restart unless-stopped \
  --pod oxpod \
  -e POSTGRESQL_ADMIN_PASSWORD=onix \
  -v oxdb:/var/lib/postgresql/data/ centos/postgresql-12-centos7

echo ==================================================
echo Starting adminer for debugging use ...
echo ==================================================
#podman run -d --name dbadmin \
#  --restart unless-stopped \
#  --pod oxpod \
#  -p 8081:8080 \ <- need to put in pod
#  adminer

echo ==================================================
echo Starting Onix database manager ...
echo ==================================================
podman run -d --name oxdbman \
  --pod oxpod \
  --rm \
  -e OX_DBM_DB_HOST=localhost \
  -e OX_DBM_DB_ADMINPWD=onix \
  -e OX_DBM_HTTP_AUTHMODE=none \
  -e OX_DBM_APPVERSION=0.0.4 \
  gatblau/dbman-snapshot

echo ==================================================
echo Initialising database with initial schema ...
echo ==================================================
podman exec -it oxdbman curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create
podman exec -it oxdbman curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy

echo ==================================================
echo Starting Onix message queue ...
echo ==================================================
podman run -d --name oxmsg \
  --restart unless-stopped \
  --pod oxpod \
  eclipse-mosquitto

echo ==================================================
echo Starting Onix ...
echo ==================================================
podman run -d --name ox \
  --restart unless-stopped \
  --pod oxpod \
  -e DB_HOST=localhost \
  -e WAPI_EVENTS_ENABLED=true \
  -e WAPI_EVENTS_SERVER_HOST=localhost \
  -e WAPI_EVENTS_SERVER_PORT=1883 \
  -e WAPI_EVENTS_SERVER_USER=admin \
  -e WAPI_EVENTS_SERVER_PWD=jdwX4HXCZGWTTD45 \
  gatblau/onix-snapshot

echo ==================================================
echo Completed startup
echo ==================================================
