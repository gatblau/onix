#!/bin/bash

. .env

# Display summary for developer/tester
echo "Technical Info"
echo
echo "Onix Swagger"
echo "http://localhost:8080/swagger-ui.html"
echo ${ONIX_HTTP_ADMIN_USER}:${ONIX_HTTP_ADMIN_PWD}
echo
echo "Pilotctl Swagger"
echo "http://localhost:8888/api/index.html"
echo ${PILOTCTL_ONIX_USER}:${PILOTCTL_ONIX_PWD}
echo
echo "Event Receiver (Mongo) Swagger"
echo "http://localhost:${PILOTCTL_EVR_MONGO_PORT}/api/index.html"
echo ${PILOTCTL_EVR_MONGO_UNAME}:${PILOTCTL_EVR_MONGO_PWD}
echo
echo "Remote Artisan Registry Links"
echo "http://artreg.apsedge.io/api/index.html"
echo "http://nexus.apsedge.io"
echo
echo "Current Container tags"
echo CIT_OX_APP -------- $( echo ${CIT_OX_APP}  | awk -F":" '{print $2}' )
echo CIT_PILOTCTL_APP -- $( echo ${CIT_PILOTCTL_APP}  | awk -F":" '{print $2}' )
echo CIT_ARTREG_APP ---- $( echo ${CIT_ARTREG_APP}  | awk -F":" '{print $2}' )
echo CIT_DBMAN --------- $( echo ${CIT_DBMAN}  | awk -F":" '{print $2}' )
echo CIT_EVRMONGO_APP -- $( echo ${CIT_EVRMONGO_APP}  | awk -F":" '{print $2}' )