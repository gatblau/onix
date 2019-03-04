FROM openjdk:8-jdk-alpine
MAINTAINER Gatblau <onix@gatblau.org>
VOLUME ./tmp
ARG DEPENDENCY=tmp
COPY ${DEPENDENCY}/BOOT-INF/lib /app/lib
COPY ${DEPENDENCY}/META-INF /app/META-INF
COPY ${DEPENDENCY}/BOOT-INF/classes /app
USER 20
ENTRYPOINT ["java","-cp","app:app/lib/*","org/gatblau/onix/App"]
