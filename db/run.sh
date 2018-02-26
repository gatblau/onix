#!/usr/bin/env bash
docker run -it -d -p 5432:5432 -e POSTGRESQL_USER=onix -e POSTGRESQL_PASSWORD=onix -e POSTGRESQL_DATABASE=onix onix-db:0.0.1-0