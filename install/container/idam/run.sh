#!/usr/bin/env bash
docker run --name idam -p 8081:8080 -d -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin jboss/keycloak