This folder contains scripts to facilitate testing of Prometheus rules and configuration of AlertManager.

| file | description |
|---|---|
| [download.sh](download.sh)| downloads prometheus, alertmanager and goreman for Mac OS. |
| [prometheus.yml](prometheus.yml) | Prometheus server configuration file. |
| [etcd_rules.yml](etcd_rules.yml) | Prometheus rules for etcd service up/down status changes. |
| [alertmanager.yml](alertmanager.yml) | AlertManager configuration file. |
| [Procfile](Procfile) | [Goreman](https://github.com/mattn/goreman) uses the Heroku Procfile to launch all required processes. |

## Setting up SeS locally

In order to set up SeS locally for configuration testing follow the steps below:

### Download required subsystems

```bash
# download prometheus, alertmanager and etcd
$ sh download.sh
```

### Run subsystems
```bash
# install goreman - needs go installed
$ go get github.com/mattn/goreman

# run Procfile launching prometheus, alertmanager and etcd cluster
$ goreman start
```


Once *etcd, prometheus and alertmanager* are downloaded run the following command to start them up:
```bash
$ goreman start
```