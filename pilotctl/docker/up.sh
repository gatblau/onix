#
# Remote Host Service Setup
# this script setup Remote Host service ready for Pilot to connect to
#
# source .env file vars
set -o allexport; source .env; set +o allexport
# if the artisan registry URI is not set then set it to the local hostname
# this is a hack to connect to a registry running in the localhost outside of the compose network
if [ -n "${PILOTCTL_ART_REG_URI+x}" ];
then
  PILOTCTL_ART_REG_URI=http://${HOSTNAME}:8082
  echo PILOTCTL_ART_REG_URI not configured, assuming art registry is listening at localhost '${PILOTCTL_ART_REG_URI}'
fi
# start all services
docker-compose up -d --remove-orphans
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

# create required test items
curl -X PUT "http://localhost:8080/item/ART_FX:LIST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/fx.json"
curl -X PUT "http://localhost:8080/item/ORG_GRP:ACME" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-grp-acme.json"
curl -X PUT "http://localhost:8080/item/ORG:OPCO_A" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-a.json"
curl -X PUT "http://localhost:8080/item/ORG:OPCO_B" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/org-opco-b.json"
# areas
curl -X PUT "http://localhost:8080/item/AREA:EAST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-east.json"
curl -X PUT "http://localhost:8080/item/AREA:WEST" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-west.json"
curl -X PUT "http://localhost:8080/item/AREA:NORTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-north.json"
curl -X PUT "http://localhost:8080/item/AREA:SOUTH" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/area-south.json"
# locations
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON:PADDINGTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-paddington.json"
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON:EUSTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-euston.json"
curl -X PUT "http://localhost:8080/item/LOCATION:LONDON:BANK" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-london-bank.json"
curl -X PUT "http://localhost:8080/item/LOCATION:MANCHESTER:PICCADILLY" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-piccadilly.json"
curl -X PUT "http://localhost:8080/item/LOCATION:MANCHESTER:CHORLTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@items/location-manchester-chorlton.json"

# create required test links
# org group -> org
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|ORG:OPCO_A" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-a.json"
curl -X PUT "http://localhost:8080/link/ORG_GRP:ACME|ORG:OPCO_B" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/acme-opco-b.json"
# org group -> area
# TODO
# org -> location
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_PADDINGTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_EUSTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:LONDON_BANK" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-a-london-paddington.json"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:MANCHESTER_PICCADILLY" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-piccadilly.json"
curl -X PUT "http://localhost:8080/link/ORG:OPCO_A|LOCATION:MANCHESTER_CHORLTON" -u "$ONIX_HTTP_ADMIN_USER:$ONIX_HTTP_ADMIN_PWD" -H  "accept: application/json" -H  "Content-Type: application/json" -d "@links/opco-b-manchester-chorlton.json"
# area -> location
# TODO
# stop dbman instances
docker-compose stop dbman_pilotctl
docker-compose stop dbman_ox