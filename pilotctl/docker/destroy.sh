#!/bin/bash

echo "===================================================="
echo "Stopping all containers ..."
docker-compose down

echo ====================================================
echo "Destroying persistent data ..."
docker volume rm ${PWD##*/}_evr-mongo-db
docker volume rm ${PWD##*/}_evr-mongo-dblogs
docker volume rm ${PWD##*/}_db
