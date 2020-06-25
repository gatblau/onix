PLUGIN_PREFIX="dbman-db-"
cd plugins/pgsql
go build -o ${PLUGIN_PREFIX}pgsql
cd ../..
mv ./plugins/pgsql/${PLUGIN_PREFIX}pgsql ./${PLUGIN_PREFIX}pgsql
go build
