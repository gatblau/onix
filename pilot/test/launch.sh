-- creates the haproxy pod
sudo podman pod create --name TCS -p 80:80,8085:8085,8080:8080

-- add a pgsql db container to the pod
sudo podman run -d --restart=always --pod=TCS -e POSTGRESQL_ADMIN_PASSWORD=onix --name=oxdb centos/postgresql-12-centos7

sleep 5

-- add a dbman container to the pod
sudo podman run -d --restart=always --pod=TCS -e OX_DBM_DB_HOST=oxdb -e OX_DBM_DB_ADMINPWD=onix -e OX_DBM_HTTP_AUTHMODE=none -e OX_DBM_APPVERSION=0.0.4 --name=dbman gatblau/dbman-snapshot

sleep 3

# deploy onix database
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/create 2>&1
curl -H "Content-Type: application/json" -X POST http://localhost:8085/db/deploy 2>&1

-- add a onix wapi container to the pod
sudo podman run \
-d --restart=always --pod=TCS \
--name=ox gatblau/onix-snapshot

-- add haproxy container to the pod
sudo podman run \
-d --restart=always --pod=TCS \
--name=haproxy haproxy

-- add a pilot container to the pod
sudo podman run \
-d --restart=always --pod=TCS \
-e OXP_ONIX_URL="http://127.0.0.1:8080" \
-e OXP_ONIX_AUTHMODE="basic" \
-e OXP_ONIX_USERNAME="admin" \
-e OXP_ONIX_PASSWORD="0n1x" \
-e OXP_BROKER_SERVER="tcp://127.0.0.1:1883" \
-e OXP_BROKER_INSECURESKIPVERIFY="true" \
-e OXP_APP_KEY="TEST_APP_01" \
-e OXP_APP_CFGFILE="/usr/local/etc/haproxy/haproxy.cfg" \
-e OXP_APP_RELOADCMD="kill -SIGHUP 1" \
--name=pilot gatblau/pilot-snapshot
