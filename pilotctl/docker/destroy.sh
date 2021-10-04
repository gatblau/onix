#!/bin/bash

docker-compose -f onix.yaml down
docker volume rm ${PWD##*/}_evr-mongo-db
docker volume rm ${PWD##*/}_ox-db
