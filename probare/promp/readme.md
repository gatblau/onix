# Prometheus Probe

A tiny programme that launches a Prometheus metrics endpoint on a configurable port.
Facilitates the testing of application availability using Prometheus.

To compile:

```bash
$ make build
```

To Run:

```bash
# launch 4 processes on 4 different ports
$ ./promp -p 8081
$ ./promp -p 8082
$ ./promp -p 8083
$ ./promp -p 8084
```