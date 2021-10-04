#!/bin/bash

docker-compose -f control-plane.yaml down
docker volume rm ${PWD##*/}_evr-mongo-db
docker volume rm ${PWD##*/}_ox-db
