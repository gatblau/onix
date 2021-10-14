#!/bin/bash

echo "===================================================="
echo "Stopping all Onix containers ..."
docker-compose down

echo "===================================================="
echo "Stopping Nexus backend ..."
docker rm nexus -f

echo "===================================================="
echo "Destroying Docker volumes ..."
docker volume rm ${PWD##*/}_evr-mongo-db
docker volume rm ${PWD##*/}_evr-mongo-dblogs
docker volume rm ${PWD##*/}_db
docker volume rm ${PWD##*/}_nexus

echo "===================================================="
echo "Destroying local flat files ..."
rm -f .env || true
rm -f conf/postgres_servers.json || true
rm -f conf/ev_received.json || true
