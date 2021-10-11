#!/bin/bash

#
#    Onix Pilot Host Control Service
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# functions
RNDPASS () {
  tr -dc A-Za-z0-9 </dev/urandom | head -c 16
}

# create new .env
echo Generating environment file ...

cat > .env <<EOF
################################################################################################
# Docker specific variables
################################################################################################
# make sure this doesn't clash with any other networks you may have to ensure data and services are sandboxed
DOCKER_NETWORK=controlplane

################################################################################################
# Container image tags
################################################################################################
CIT_MONGO=mongo:5
CIT_POSTGRES=postgres:13
CIT_OX_APP=quay.io/gatblau/onix-snapshot:v0.0.4-1af14bb-021021131813
CIT_PILOTCTL_APP=quay.io/gatblau/pilotctl:0.0.4-081021093913126-9ea4c9e2bd
CIT_ARTREG_APP=quay.io/gatblau/artisan-registry:0.0.4-011021162133879-a3dedecb3f-RC1
CIT_DBMAN=quay.io/gatblau/dbman-snapshot:v0.0.4-d4fb6f7-031020001129
CIT_EVRMONGO_APP=quay.io/gatblau/pilotctl-evr-mongodb:0.0.4-300921174051295-11aab8b6cc

################################################################################################
# Postgres (used for Onix and Pilotctl)
################################################################################################

PG_ADMIN_USER=postgres
PG_ADMIN_PWD=$(RNDPASS)

################################################################################################
# DBMan - (@ localhost:8085/api/)
################################################################################################

# Credentials
DBMAN_HTTP_USER=admin
DBMAN_HTTP_PWD=$(RNDPASS)

# the authentication mode used by dbman (e.g. none or basic)
DBMAN_AUTH_MODE=basic

# the onix application version dbman uses to retrieve the database scripts
# latest is 0.0.4 which corresponds to database schema 4
# https://github.com/gatblau/ox-db/blob/master/plan.json
DBMAN_ONIX_VERSION=0.0.4
DBMAN_PILOTCTL_VERSION=0.0.4

# Other
DBMAN_PILOTCTL_COMMIT_HASH=master # master is the latest version, enter hash if different is required
DBMAN_PILOTCTL_REPO_URI=https://raw.githubusercontent.com/gatblau/pilotctl-db/
DBMAN_ONIX_COMMIT_HASH=master # master is the latest version, enter hash if different is required
DBMAN_ONIX_REPO_URI=https://raw.githubusercontent.com/gatblau/ox-db/

################################################################################################
# Onix - (@ localhost:8080/swagger-ui.html)
################################################################################################

# Credentials
# Notes: Do not changes user names but recommend change passwords
ONIX_DB_USER=onix
ONIX_DB_PWD=$(RNDPASS)
ONIX_HTTP_ADMIN_USER=admin
ONIX_HTTP_ADMIN_PWD=$(RNDPASS)
ONIX_HTTP_READER_USER=reader
ONIX_HTTP_READER_PWD=$(RNDPASS)
ONIX_HTTP_WRITER_USER=writer
ONIX_HTTP_WRITER_PWD=$(RNDPASS)

# Other
AUTH_MODE=basic # the authentication mode used by the Onix Web API (set to Basic Authentication)

################################################################################################
# Pilotctl
################################################################################################

# DB Credentials
# Notes: Do not changes user names but recommend change passwords
PILOTCTL_DB_USER=pilotctl
PILOTCTL_DB_PWD=$(RNDPASS)

# Service credentials
PILOTCTL_HTTP_USER=pilotctl
PILOTCTL_HTTP_PWD=$(RNDPASS)
PILOTCTL_HTTP_PORT=8888

# Swagger (@ localhost:8888/api/index.html)
# NB. Temporary creds until RBAC version has completed testing & released
PILOTCTL_ONIX_USER=admin@pilotctl.com # used for authentication - could be different than email if required
PILOTCTL_ONIX_EMAIL=admin@pilotctl.com # used for password resets
PILOTCTL_ONIX_PWD=P1l0tctl

################################################################################################
# Artisan Registry
################################################################################################
ART_REG_URI=https://artreg.apsedge.io
ART_REG_USER=admin
ART_REG_PWD=$(RNDPASS)
ART_REG_PORT=443
#ART_REG_BACKEND_URI=https://nexus.apsedge.io
#ART_REG_BACKEND_PORT=443

################################################################################################
# PILOTCTL backing artisan registry
################################################################################################
#PILOTCTL_ART_REG_URI=http://artreg-app:8082
#PILOTCTL_ART_REG_USER=admin
#PILOTCTL_ART_REG_PWD=$(RNDPASS)

################################################################################################
# PILOTCTL Events Receiver (Mongo version)
################################################################################################
PILOTCTL_EVR_MONGO_APPCONTAINER=evr-mongo-app
PILOTCTL_EVR_MONGO_DBCONTAINER=evr-mongo-db
PILOTCTL_EVR_MONGO_OPTIONS=/syslog?authSource=admin&keepAlive=true&poolSize=30&autoReconnect=true&socketTimeoutMS=360000&connectTimeoutMS=360000
PILOTCTL_EVR_MONGO_DBPORT=27017
PILOTCTL_EVR_MONGO_PORT=8885
PILOTCTL_EVR_MONGO_UNAME=admin
PILOTCTL_EVR_MONGO_PWD=$(RNDPASS)


################################################################################################
# General
################################################################################################

LOGIN_LEVEL=Trace

# the URI of the Onix Web API
WAPI_URL=http://ox-app
WAPI_PORT=8080

# enables publication of events to the MQTT message broker when items of specified type change
BROKER_ENABLED=false
BROKER_PORT=1883
EOF

# create new docker-file
echo Generating docker-compose file ...
cat > docker-compose.yaml <<EOF
#
#    Onix Pilot Host Control Service
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# Note: Some ports in this Compose are deliberately exposed to the host where they
#       wouldn't normally be for a "production" system. This is so that developers
#       or technically minded users can test or debug areas of the system.

version: '3'

services:

  #############################################################################
  # Onix application services
  #############################################################################
  ox-app:
    image: "${CIT_OX_APP}"
    depends_on:
      - db
      - ox-dbman
    container_name: ox-app
    restart: always
    environment:
      - DB_HOST=db
      - DB_USER=${ONIX_DB_USER}
      - DB_PWD=${ONIX_DB_PWD}
      - DB_ADMIN_USER=${PG_ADMIN_USER}
      - DB_ADMIN_PWD=${PG_ADMIN_PWD}
      - WAPI_AUTH_MODE=${AUTH_MODE}
      - WAPI_ADMIN_USER=${ONIX_HTTP_ADMIN_USER}
      - WAPI_ADMIN_PWD=${ONIX_HTTP_ADMIN_PWD}
      - WAPI_EVENTS_ENABLED=${BROKER_ENABLED}
      - WAPI_EVENTS_SERVER_HOST=oxmsg
      - WAPI_EVENTS_SERVER_PORT=${BROKER_PORT}
    ports:
      - "8080:8080"

  pilotctl-app:
    image: ${CIT_PILOTCTL_APP}
    depends_on:
      - db
      - ox-app
      - pilotctl-dbman
    container_name: pilotctl-app
    restart: always
    environment:
      - OX_PILOTCTL_DB_HOST=db
      - OX_PILOTCTL_DB_USER=${PILOTCTL_DB_USER}
      - OX_PILOTCTL_DB_PWD=${PILOTCTL_DB_PWD}
      - OX_HTTP_UNAME=${PILOTCTL_HTTP_USER}
      - OX_HTTP_PWD=${PILOTCTL_HTTP_PWD}
      - OX_HTTP_PORT=${PILOTCTL_HTTP_PORT}
      - OX_WAPI_URI=${WAPI_URL}:${WAPI_PORT}
      - OX_WAPI_USER=${ONIX_HTTP_ADMIN_USER}
      - OX_WAPI_PWD=${ONIX_HTTP_ADMIN_PWD}
      - OX_WAPI_INSECURE_SKIP_VERIFY=true
      - OX_ART_REG_URI=${ART_REG_URI}
      - OX_ART_REG_USER=${ART_REG_USER}
      - OX_ART_REG_PWD=${ART_REG_PWD}
    ports:
      - "8888:8888"
    volumes:
      - ./keys:/keys
      - ./conf:/conf

  # artreg-app:
  #   image: ${CIT_ARTREG_APP}
  #   container_name: artreg-app
  #   restart: always
  #   environment:
  #     - OXA_HTTP_UNAME=${ART_REG_USER}
  #     - OXA_HTTP_PORT=${ART_REG_PORT}
  #     - OXA_HTTP_PWD=${ART_REG_PWD}
  #     - OXA_HTTP_BACKEND_DOMAIN=${ART_REG_BACKEND_URI}:${ART_REG_BACKEND_PORT}
  #   ports:
  #     - 8082:8082

  #############################################################################
  # Temporary utility services
  #############################################################################
  ox-dbman:
    image: ${CIT_DBMAN}
    container_name: ox-dbman
    restart: always
    environment:
      - OX_DBM_DB_HOST=db
      - OX_DBM_DB_USERNAME=${ONIX_DB_USER}
      - OX_DBM_DB_PASSWORD=${ONIX_DB_PWD}
      - OX_DBM_DB_ADMINUSERNAME=${PG_ADMIN_USER}
      - OX_DBM_DB_ADMINPASSWORD=${PG_ADMIN_PWD}
      - OX_DBM_HTTP_USERNAME=${DBMAN_HTTP_USER}
      - OX_DBM_HTTP_PASSWORD=${DBMAN_HTTP_PWD}
      - OX_DBM_HTTP_AUTHMODE=${DBMAN_AUTH_MODE}
      - OX_DBM_APPVERSION=${DBMAN_ONIX_VERSION}
      - OX_DBM_REPO_URI=${DBMAN_ONIX_REPO_URI}${DBMAN_ONIX_COMMIT_HASH}
    ports:
      - "8085:8085"

  pilotctl-dbman:
    image: ${CIT_DBMAN}
    container_name: pilotctl-dbman
    restart: always
    environment:
      - OX_DBM_DB_HOST=db
      - OX_DBM_DB_NAME=pilotctl
      - OX_DBM_DB_USERNAME=${PILOTCTL_DB_USER}
      - OX_DBM_DB_PASSWORD=${PILOTCTL_DB_PWD}
      - OX_DBM_DB_ADMINUSERNAME=${PG_ADMIN_USER}
      - OX_DBM_DB_ADMINPASSWORD=${PG_ADMIN_PWD}
      - OX_DBM_HTTP_USERNAME=${DBMAN_HTTP_USER}
      - OX_DBM_HTTP_PASSWORD=${DBMAN_HTTP_PWD}
      - OX_DBM_HTTP_AUTHMODE=${DBMAN_AUTH_MODE}
      - OX_DBM_APPVERSION=${DBMAN_PILOTCTL_VERSION}
      - OX_DBM_REPO_URI=${DBMAN_PILOTCTL_REPO_URI}${DBMAN_PILOTCTL_COMMIT_HASH}
    ports:
      - "8086:8085"

  #############################################################################
  # Database services
  #############################################################################
  db: # (supports Onix and Pilot Control)
    image: ${CIT_POSTGRES}
    container_name: db
    restart: always
    environment:
      - POSTGRES_PASSWORD=${PG_ADMIN_PWD}
    volumes:
      - db:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  #############################################################################
  # Event Receiver services
  #############################################################################
  evr-mongo-app:
    image: ${CIT_EVRMONGO_APP}
    depends_on:
      - pilotctl-app
      - evr-mongo-db
    container_name: evr-mongo-app
    restart: always
    environment:
      - OX_MONGO_EVR_CONN=mongodb://${PILOTCTL_EVR_MONGO_UNAME}:${PILOTCTL_EVR_MONGO_PWD}@${PILOTCTL_EVR_MONGO_DBCONTAINER}:${PILOTCTL_EVR_MONGO_DBPORT}/${PILOTCTL_EVR_MONGO_OPTIONS}
      - OX_HTTP_PORT=${PILOTCTL_EVR_MONGO_PORT}
      - OX_HTTP_UNAME=${PILOTCTL_EVR_MONGO_UNAME}
      - OX_HTTP_PWD=${PILOTCTL_EVR_MONGO_PWD}
    ports:
      - "${PILOTCTL_EVR_MONGO_PORT}:8885"

  evr-mongo-db:
    image: ${CIT_MONGO}
    container_name: evr-mongo-db
    restart: always
    environment:
      - MONGO_INITDB_ROOT_USERNAME=${PILOTCTL_EVR_MONGO_UNAME}
      - MONGO_INITDB_ROOT_PASSWORD=${PILOTCTL_EVR_MONGO_PWD}
    ports:
      - ${PILOTCTL_EVR_MONGO_DBPORT}:27017
    volumes:
      - evr-mongo-db:/data/db
      - evr-mongo-dblogs:/var/log/mongodb


#############################################################################
# Networking
#############################################################################
networks:
  default:
    external:
      name: ${DOCKER_NETWORK}

#############################################################################
# Data volumes
#############################################################################
volumes:
  db:
  evr-mongo-db:
  evr-mongo-dblogs:
EOF

# source the newly created .env file
set -o allexport; source .env; set +o allexport

# create events receiver JSON for all events receivers
[ ! -d "./conf" ] && mkdir conf
cat > ./conf/ev_receive.json <<EOF
{
  "event_receivers": [
    {
      "uri": "http://${PILOTCTL_EVR_MONGO_APPCONTAINER}:${PILOTCTL_EVR_MONGO_PORT}/events",
      "user": "${PILOTCTL_EVR_MONGO_UNAME}",
      "pwd": "${PILOTCTL_EVR_MONGO_PWD}"
    }
  ]
}
EOF

# Ensure attachable Docker network is already created
if [[ $(docker network inspect ${DOCKER_NETWORK}) == "[]" ]]; then
  echo Creating Docker network ${DOCKER_NETWORK} ...
  docker network create ${DOCKER_NETWORK}
fi

# start all services
docker-compose up -d

# NB. If you have issues with either Onix admin password not being updated or pilotctl user not being created,
# try increasing timeout here to ensure that service has time to spin up correctly before applying data changes
# (The time taken will depend on your host computing and I/O)
echo Waiting 5s for configuration 1 to apply ...
sleep 5

# setup the onix database
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1
# setup the rem database
curl -H "Content-Type: application/json" -X POST http://localhost:8086/db/create 2>&1
curl -H "Content-Type: application/json" -X POST http://localhost:8086/db/deploy 2>&1

echo Waiting 20s for configuration 2 to apply ...
sleep 20

# update default's Onix Web API admin password"
curl  --connect-timeout 5 \
      --max-time 10 \
      --retry 5 \
      --retry-delay 0 \
      --retry-max-time 30 \
      -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" \
      -H "Content-Type: application/json" \
      -X PUT http://localhost:8080/user/$ONIX_HTTP_ADMIN_USER/pwd \
      -d "{\"pwd\":\"$ONIX_HTTP_ADMIN_PWD\"}"

# create pilotctl user
curl  --connect-timeout 5 \
      --max-time 10 \
      --retry 5 \
      --retry-delay 0 \
      --retry-max-time 30 \
      -H "Authorization: Basic $(printf '%s:%s' $ONIX_HTTP_ADMIN_USER $ONIX_HTTP_ADMIN_PWD | base64)" \
      -H "Content-Type: application/json" \
      -X PUT http://localhost:8080/user/ONIX_PILOTCTL \
      -d "{\"email\":\"${PILOTCTL_ONIX_EMAIL}\", \"name\":\"${PILOTCTL_ONIX_USER}\", \"pwd\":\"${PILOTCTL_ONIX_PWD}\", \"service\":\"false\", \"acl\":\"*:*:*\"}"

# create required test items
curl -X PUT "http://localhost:8080/item/ART_FX:LIST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/fx.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/ORG_GRP:ACME" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-grp-acme.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/ORG:OPCO_A" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-a.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/ORG:OPCO_B" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-b.json" && printf "\n"
# areas
curl -X PUT "http://localhost:8080/item/AREA:EAST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-east.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/AREA:WEST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-west.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/AREA:NORTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-north.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/AREA:SOUTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-south.json" && printf "\n"
# locations
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON_PADDINGTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-paddington.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON_EUSTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-euston.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON_BANK" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-bank.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/LOCATION:MANCHESTER_PICCADILLY" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-piccadilly.json" && printf "\n"
curl -X PUT "http://localhost:8080/item/LOCATION:MANCHESTER_CHORLTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-chorlton.json" && printf "\n"

# create required test links
# org group -> org
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|ORG:OPCO_A" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-a.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|ORG:OPCO_B" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-b.json" && printf "\n"
# org group -> area
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|AREA:EAST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-east.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|AREA:WEST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-west.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|AREA:NORTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-north.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|AREA:SOUTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-south.json" && printf "\n"
# org -> location
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_PADDINGTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_EUSTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_BANK" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:MANCHESTER_PICCADILLY" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-piccadilly.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:MANCHESTER_CHORLTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-chorlton.json" && printf "\n"
# area -> location
curl -X PUT "http://localhost:8080/link/AREA:SOUTH|LOCATION:LONDON_PADDINGTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-paddington.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/AREA:SOUTH|LOCATION:LONDON_EUSTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-euston.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/AREA:SOUTH|LOCATION:LONDON_BANK" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-bank.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/AREA:NORTH|LOCATION:MANCHESTER_PICCADILLY" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/north-manchester-piccadilly.json" && printf "\n"
curl -X PUT "http://localhost:8080/link/AREA:NORTH|LOCATION:MANCHESTER_CHORLTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/north-manchester-chorlton.json" && printf "\n"

# stop dbman instances
docker-compose stop pilotctl-dbman
docker-compose stop ox-dbman

# Completed
echo ====================================================================================
echo Deploy is completed
