#!/bin/bash

. .env

# Display summary for developer/tester
echo "Developer/Tester info"
echo
echo "Onix (Swagger)"
echo "http://localhost:8080/swagger-ui.html"
echo ${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}
echo
echo "Pilotctl (Swagger)"
echo "http://localhost:8888/api/index.html"
echo ${PILOTCTL_ONIX_USER}:${PILOTCTL_ONIX_PWD}
