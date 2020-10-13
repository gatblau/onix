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
$ curl localhost:8081
$ curl localhost:8082
$ curl localhost:8083
$ curl localhost:8084
```