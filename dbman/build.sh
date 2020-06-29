#!/bin/sh
PLUGIN_PREFIX="dbman-db-"
cd plugin/pgsql
go build -o ${PLUGIN_PREFIX}pgsql
cd ../..
mv ./plugin/pgsql/${PLUGIN_PREFIX}pgsql ./${PLUGIN_PREFIX}pgsql
go build .
