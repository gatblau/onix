# Prometheus Probe

A tiny programme that launches a Prometheus metrics endpoint on a configurable port.
Facilitates the testing of application availability using Prometheus.

## Compile:

```bash
$ make build
```

## Run

```bash
# launch 4 processes on 4 different ports
$ ./promp -p 8081
$ ./promp -p 8082
$ ./promp -p 8083
$ ./promp -p 8084
```

## Access the Prometheus metrics endpoint

```bash
$ curl localhost:8081/metrics
$ curl localhost:8082/metrics
$ curl localhost:8083/metrics
$ curl localhost:8084/metrics
```

## As a docker image

```bash
$ docker run -it --rm --name test01 -p 3000:8081 gatblau/promp-snapshot
$ docker run -it --rm --name test02 -p 3000:8082 gatblau/promp-snapshot
$ docker run -it --rm --name test03 -p 3000:8083 gatblau/promp-snapshot
$ docker run -it --rm --name test04 -p 3000:8084 gatblau/promp-snapshot
```