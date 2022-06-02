#!/usr/bin/env bash

# launch mongo db
docker run \
  --name mongo \
  -p 27017:27017 \
  -e "MONGO_INITDB_DATABASE=doorman" \
  -e "MONGO_INITDB_ROOT_USERNAME=admin" \
  -e "MONGO_INITDB_ROOT_PASSWORD=admin" \
  -d \
  mongo

# launch ultralight registry
docker run \
    --name uar \
    -p 8082:8080 \
    -e "OX_ADMIN_USER=admin" \
    -e "OX_ADMIN_PWD=admin" \
    -d \
    quay.io/artisan/uar

# launch container registry
docker run \
  --name registry \
  -p 5000:5000 \
  --restart always \
  -d \
  registry:2