#!/bin/bash

OPTIONS='%-40s | %-40s | %-40s \n'

. .env

# Display summary for developer/tester
echo "Credential Info"
echo

printf "$OPTIONS" "------------" "------------" "------------"
printf "$OPTIONS" "Service" "URI" "Credentials"
printf "$OPTIONS" "------------" "------------" "------------"
printf "$OPTIONS" "Onix Swagger" "http://localhost:8080/swagger-ui.html" "${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}"
printf "$OPTIONS" "Pilotctl Swagger" "http://localhost:8888/api/index.html" "${PILOTCTL_ONIX_USER}:${PILOTCTL_ONIX_PWD}"
printf "$OPTIONS" "Event Receiver (Mongo) Swagger" "http://localhost:${PILOTCTL_EVR_MONGO_PORT}/api/index.html" "${PILOTCTL_EVR_MONGO_UNAME}:${PILOTCTL_EVR_MONGO_PWD}"
printf "$OPTIONS" "Local Artisan Swagger" "http://localhost:${ART_REG_PORT}/api/index.html" "${ART_REG_USER}:${ART_REG_PWD}"
printf "$OPTIONS" "PilotCtl database (direct Postgres)" "http://localhost:5432" "${PG_ADMIN_USER}:${PG_ADMIN_PWD}"
printf "$OPTIONS" "Pilotctl database (Web GUI)" "http://localhost:8083" "admin@local.com:${PG_ADMIN_PWD}"
printf "$OPTIONS" "------------" "------------" "------------"
printf "$OPTIONS" "Demo Artisan Registry" "artisan.apsedge.io" "Please contact for credentials"
printf "$OPTIONS" "------------" "------------" "------------"

echo
echo "Current Container tags"
cat .env | grep  --color=never CIT_ | sort
