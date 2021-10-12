#!/bin/bash

#
#    Onix Pilot Host Control Service
#    Copyright (c) 2018-2021 by www.gatblau.org
#    Licensed under the Apache License, Version 2.0 at http://www.apache.org/licenses/LICENSE-2.0
#    Contributors to this project, hereby assign copyright in this code to the project,
#    to be licensed under the same terms as the rest of the code.
#

# CURL options
OPTIONS="-v --silent --show-error -f -L --max-redirs 3 --retry 10 --retry-connrefused --retry-delay 5 --max-time 120"
CURL_MAXRETRY=10
CURL_DELAY=5


# functions
RNDPASS () {
  tr -dc A-Za-z0-9 </dev/urandom | head -c 16
}

CURL2xx() {
  CURL_CODE=0
  CURL_RETRY=0
  while [[ ("$CURL_CODE" < 200 || "$CURL_CODE" > 299) && "$CURL_RETRY" != "$CURL_MAXRETRY" ]]
  do
    #echo DBG curl -s -o /dev/nul -w %{http_code} "$@"
    CURL_CODE=$(curl -s -o /dev/nul -w %{http_code} "$@")
    if [[ ("$CURL_CODE" < 200 || "$CURL_CODE" > 299) && "$CURL_RETRY" != "$CURL_MAXRETRY" ]]
    then
      ((CURL_RETRY=CURL_RETRY+1))
      echo "Got return code $CURL_CODE, pausing for retry $CURL_RETRY ..."
      sleep $CURL_DELAY
    fi
  done
}


# create new .env
echo Generating environment file

cat > .env <<EOF
################################################################################################
# Please check/amend these details below
################################################################################################

# Artisan registry to use
ART_REG_URI=http://artreg-app
ART_REG_USER=admin
ART_REG_PWD=$(RNDPASS)
ART_REG_PORT=8082
ART_REG_BACKEND_URI=http://nexus
ART_REG_BACKEND_PORT=8081


# Docker network to use
DOCKER_NETWORK=onix

# Container image tags
CIT_MONGO=docker.io/mongo:5
CIT_MONGOGUI=docker.io/mongo-express:latest
CIT_POSTGRES=docker.io/postgres:13
CIT_POSTGRESGUI=docker.io/dpage/pgadmin4:latest
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
PILOTCTL_HTTP_PORT=8888

# NB. Temporary creds until RBAC version has completed testing & released
PILOTCTL_ONIX_URI=http://ox-app:8080
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
WAPI_URL=http://localhost
WAPI_PORT=8080
EOF

echo "Environment file (.env) has been created"

# source the automatically created .env file
set -o allexport; source .env; set +o allexport

# create JSON file for PGAdmin GUI
cat > ./conf/postgres_servers.json <<EOF
{
    "Servers": {
        "1": {
            "Name": "Main Onix database",
            "Group": "Onix",
            "Port": 5432,
            "Username": "postgres",
            "Host": "db",
            "MaintenanceDB": "postgres",
            "Username": "postgres",
            "SSLMode": "disable"
        }
    }
}
EOF

# Ensure attachable Docker network is already created
if [[ $(docker network inspect ${DOCKER_NETWORK}) == "[]" ]]; then
  echo Creating Docker network ${DOCKER_NETWORK} ...
  docker network create ${DOCKER_NETWORK}
fi

# create Nexus backend (for Artisan Registry)
docker run -d \
  --rm \
  -p ${ART_REG_BACKEND_PORT}:8081 \
  --name nexus \
  -v ${PWD##*/}_nexus:/nexus-data \
  --network ${DOCKER_NETWORK} \
  sonatype/nexus3
echo "Waiting for API interface to be available in Nexus container ..."
sleep 15
CURRENTPASS=
echo "Checking for temporary password file ..."
while [ -z "$CURRENTPASS" ]
do
  sleep 2
  CURRENTPASS=$(docker exec nexus cat /nexus-data/admin.password)
done

echo "Updating admin password ..."
CURL2xx -X PUT \
  -u admin:${CURRENTPASS} \
  http://localhost:${ART_REG_BACKEND_PORT}/service/rest/v1/security/users/admin/change-password \
  -H 'accept: application/json' \
  -H 'Content-Type: text/plain' \
  -d "${ART_REG_PWD}"

echo "Creating new Artisan repository ..."
CURL2xx -X POST \
  -u admin:${ART_REG_PWD} \
  http://localhost:${ART_REG_BACKEND_PORT}/service/rest/v1/repositories/raw/hosted \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "name": "artisan",
  "online": true,
  "storage": {
    "blobStoreName": "default",
    "strictContentTypeValidation": true,
    "writePolicy": "allow_once"
  },
  "cleanup": {
    "policyNames": [
      "string"
    ]
  },
  "component": {
    "proprietaryComponents": true
  },
  "raw": {
    "contentDisposition": "ATTACHMENT"
  }
}'
echo "Disabling anonymous access ..."
CURL2xx -X PUT \
  -u admin:${ART_REG_PWD} \
  http://localhost:${ART_REG_BACKEND_PORT}/service/rest/v1/security/anonymous \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{"enabled": false}'

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

# start all other services
docker-compose up -d

# setup the onix database
echo Creating Onix database via DBMan ...
curl $OPTIONS \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8085/db/create 2>&1
curl $OPTIONS \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8085/db/deploy 2>&1

echo Creating Pilotctl database via DBMan ...
curl $OPTIONS \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8086/db/create 2>&1
curl $OPTIONS \
  -H "Content-Type: application/json" \
  -X POST http://localhost:8086/db/deploy 2>&1


echo "Updating Onix admin password from default ..."
until contents=$(curl \
  -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" \
  -H "Content-Type: application/json" \
  -X PUT ${WAPI_URL}:${WAPI_PORT}/user/${ONIX_HTTP_ADMIN_USER}/pwd \
  -d "{\"pwd\":\"${ONIX_HTTP_ADMIN_PWD}\"}"
)
do
  sleep 3
done

echo "Creating special pilotctl user in Onix ..."
until contents=$(curl \
  -H "Authorization: Basic $(printf '%s:%s' ${ONIX_HTTP_ADMIN_USER} ${ONIX_HTTP_ADMIN_PWD} | base64)" \
  -H "Content-Type: application/json" \
  -X PUT ${WAPI_URL}:${WAPI_PORT}/user/ONIX_PILOTCTL \
  -d "{\"email\":\"${PILOTCTL_ONIX_EMAIL}\", \"name\":\"${PILOTCTL_ONIX_USER}\", \"pwd\":\"${PILOTCTL_ONIX_PWD}\", \"service\":\"false\", \"acl\":\"*:*:*\"}"
)
do
  sleep 3
done

# create required test items
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/ART_FX:LIST" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/fx.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/ORG_GRP:ACME" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-grp-acme.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/ORG:OPCO_A" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-a.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/ORG:OPCO_B" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-b.json" && printf "\n"
# areas
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/AREA:EAST" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-east.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/AREA:WEST" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-west.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/AREA:NORTH" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-north.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/AREA:SOUTH" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-south.json" && printf "\n"
# locations
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/LOCATION:LONDON_PADDINGTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-paddington.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/LOCATION:LONDON_EUSTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-euston.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/LOCATION:LONDON_BANK" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-bank.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/LOCATION:MANCHESTER_PICCADILLY" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-piccadilly.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/item/LOCATION:MANCHESTER_CHORLTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-chorlton.json" && printf "\n"

# create required test links
# org group -> org
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|ORG:OPCO_A" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-a.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|ORG:OPCO_B" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-b.json" && printf "\n"
# org group -> area
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|AREA:EAST" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-east.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|AREA:WEST" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-west.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|AREA:NORTH" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-north.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG_GRP:ACME|AREA:SOUTH" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-south.json" && printf "\n"
# org -> location
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG:OPCO_A|LOCATION:LONDON_PADDINGTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG:OPCO_A|LOCATION:LONDON_EUSTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG:OPCO_A|LOCATION:LONDON_BANK" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG:OPCO_A|LOCATION:MANCHESTER_PICCADILLY" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-piccadilly.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/ORG:OPCO_A|LOCATION:MANCHESTER_CHORLTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-chorlton.json" && printf "\n"
# area -> location
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/AREA:SOUTH|LOCATION:LONDON_PADDINGTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-paddington.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/AREA:SOUTH|LOCATION:LONDON_EUSTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-euston.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/AREA:SOUTH|LOCATION:LONDON_BANK" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/south-london-bank.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/AREA:NORTH|LOCATION:MANCHESTER_PICCADILLY" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/north-manchester-piccadilly.json" && printf "\n"
curl $OPTIONS -X PUT "${WAPI_URL}:${WAPI_PORT}/link/AREA:NORTH|LOCATION:MANCHESTER_CHORLTON" -u "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/north-manchester-chorlton.json" && printf "\n"

# stop dbman instances
docker-compose stop pilotctl-dbman
docker-compose stop ox-dbman

# Display help
echo "The deployment is now complete. Please run the info.sh helper (E.G. sh ./info.sh) script to show the credentials that have been generated for your local instances."
