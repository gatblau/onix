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
printf "$OPTIONS" "------------" "------------" "------------"
printf "$OPTIONS" "Main databases (direct Postgres)" "http://localhost:5432" "${PG_ADMIN_USER}:${PG_ADMIN_PWD}"
printf "$OPTIONS" "Main databases (Web GUI)" "http://localhost:8083" "admin@local.com:${PG_ADMIN_PWD}"
printf "$OPTIONS" "------------" "------------" "------------"
echo "Notes:"
echo
echo "- If you are browsing to your containers running on a remote Virtual Machine as opposed to locally, don't forget to replace localhost with the name or IP address of your Virtual Machine"
echo "- To access the main databases via the web GUI above, create a server entry in the Web GUI to a Postgresql server called db using the direct Postgres credentials shown"
echo "- A testing Artisan Registry is available, please contact for URL and credentials if you would prefer to use that rather than the local registry created"
echo
echo "Current Container tags"
cat .env | grep  --color=never CIT_ | sort
