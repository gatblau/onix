#
# Remote Host Service Setup
# this script setup Remote Host service ready for Pilot to connect to
#
# source .env file vars
set -o allexport; source .env; set +o allexport
# start all services
docker-compose up -d
# setup the onix database
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1
# setup the rem database
curl -H "Content-Type: application/json" -X POST http://localhost:8086/db/create 2>&1
curl -H "Content-Type: application/json" -X POST http://localhost:8086/db/deploy 2>&1
# updates default's Onix Web API admin password"
until contents=$(curl -H "Authorization: Basic $(printf '%s:%s' admin 0n1x | base64)" -H "Content-Type: application/json" -X PUT http://localhost:8080/user/$ONIX_HTTP_ADMIN_USER/pwd -d "{\"pwd\":\"$ONIX_HTTP_ADMIN_PWD\"}")
do
  sleep 3
done
# stop dbman instances
docker-compose stop dbman_rem
docker-compose stop dbman_ox