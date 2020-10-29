cd ../..
make build-linux
mv artie images/java-11
cd images/java-11
docker build -t gatblau/artie:java-11 .