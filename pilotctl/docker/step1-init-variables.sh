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
echo Generating environment file

cat > .env <<EOF
################################################################################################
# Please check/amend these details below
################################################################################################

# Artisan registry to use
ART_REG_URI=
ART_REG_USER=admin
ART_REG_PWD=
ART_REG_PORT=443
#PILOTCTL_ART_REG_URI=http://artreg-app:8082
#PILOTCTL_ART_REG_USER=admin
#PILOTCTL_ART_REG_PWD=$(RNDPASS)

# Docker network to use
DOCKER_NETWORK=onix

# Container image tags
CIT_MONGO=mongo:5
CIT_POSTGRES=postgres:13
CIT_OX_APP=quay.io/gatblau/onix-snapshot:v0.0.4-1af14bb-021021131813
CIT_PILOTCTL_APP=quay.io/gatblau/pilotctl:0.0.4-081021093913126-9ea4c9e2bd
CIT_ARTREG_APP=quay.io/gatblau/artisan-registry:0.0.4-011021162133879-a3dedecb3f-RC1
CIT_DBMAN=quay.io/gatblau/dbman-snapshot:v0.0.4-d4fb6f7-031020001129
CIT_EVRMONGO_APP=quay.io/gatblau/pilotctl-evr-mongodb:0.0.4-300921174051295-11aab8b6cc

################################################################################################
# Everything below this point should not normally need to be amended
################################################################################################

# Postgres (used for Onix and Pilotctl)
PG_ADMIN_USER=postgres
PG_ADMIN_PWD=$(RNDPASS)

# DBMan - (@ localhost:8085/api/)
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
DBMAN_PILOTCTL_COMMIT_HASH=master # master is the latest version, enter hash if different is required
DBMAN_PILOTCTL_REPO_URI=https://raw.githubusercontent.com/gatblau/pilotctl-db/
DBMAN_ONIX_COMMIT_HASH=master # master is the latest version, enter hash if different is required
DBMAN_ONIX_REPO_URI=https://raw.githubusercontent.com/gatblau/ox-db/

# Onix - (@ localhost:8080/swagger-ui.html)
ONIX_DB_USER=onix
ONIX_DB_PWD=$(RNDPASS)
ONIX_HTTP_ADMIN_USER=admin
ONIX_HTTP_ADMIN_PWD=$(RNDPASS)
ONIX_HTTP_READER_USER=reader
ONIX_HTTP_READER_PWD=$(RNDPASS)
ONIX_HTTP_WRITER_USER=writer
ONIX_HTTP_WRITER_PWD=$(RNDPASS)
AUTH_MODE=basic # the authentication mode used by the Onix Web API (set to Basic Authentication)

# Pilotctl
PILOTCTL_DB_USER=pilotctl
PILOTCTL_DB_PWD=$(RNDPASS)
PILOTCTL_HTTP_USER=pilotctl
PILOTCTL_HTTP_PWD=$(RNDPASS)
PILOTCTL_HTTP_PORT=8888
# NB. Temporary creds until RBAC version has completed testing & released
PILOTCTL_ONIX_USER=admin@pilotctl.com # used for authentication - could be different than email if required
PILOTCTL_ONIX_EMAIL=admin@pilotctl.com # used for password resets
PILOTCTL_ONIX_PWD=P1l0tctl

# PILOTCTL Events Receiver (Mongo version)
PILOTCTL_EVR_MONGO_APPCONTAINER=evr-mongo-app
PILOTCTL_EVR_MONGO_DBCONTAINER=evr-mongo-db
PILOTCTL_EVR_MONGO_OPTIONS=/syslog?authSource=admin&keepAlive=true&poolSize=30&autoReconnect=true&socketTimeoutMS=360000&connectTimeoutMS=360000
PILOTCTL_EVR_MONGO_DBPORT=27017
PILOTCTL_EVR_MONGO_PORT=8885
PILOTCTL_EVR_MONGO_UNAME=admin
PILOTCTL_EVR_MONGO_PWD=$(RNDPASS)

# MQTT message broker
# NB. enables publication of events to the MQTT message broker when items of specified type change
BROKER_ENABLED=false
BROKER_PORT=1883

# General
LOGIN_LEVEL=Trace
WAPI_URL=http://ox-app
WAPI_PORT=8080
EOF

echo "Environment file (.env) has been created"
echo "IMPOPRTANT --> Now go ahead and check/amend accordingly before running deploy script"