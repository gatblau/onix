docker rm -f nexus
docker run -d -p 8081:8081 --name nexus sonatype/nexus3
echo please wait for nexus to start up
sleep 60
docker container exec nexus cat ./nexus-data/admin.password