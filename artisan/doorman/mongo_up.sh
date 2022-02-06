# starts a mongo db container for development purposes
# use below connection string:
# => mongodb://admin:adm1n@127.0.0.1:27017/syslog?keepAlive=true&poolSize=30&autoReconnect=true&socketTimeoutMS=360000&connectTimeoutMS=360000
# see https://hub.docker.com/_/mongo/ for docker image documentation
docker run --name mongo-evr -d -p 27017:27017 \
  -e MONGO_INITDB_DATABASE=doorman \
  mongo

#  -e MONGO_INITDB_ROOT_USERNAME=admin \
#  -e MONGO_INITDB_ROOT_PASSWORD=adm1n \