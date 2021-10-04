#!/bin/bash

. .env

# Display summary for developer/tester
echo "Developer/Tester info"
echo
echo "Onix (Swagger)"
echo "http://localhost:8080/swagger-ui.html"
echo ${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}
echo
echo "Control Plane (Swagger / Dashboard)"
echo "- API http://localhost:8888/api/index.html"
echo "- Dashboard http://localhost"
echo ${PILOTCTL_ONIX_USER}:${PILOTCTL_ONIX_PWD}
