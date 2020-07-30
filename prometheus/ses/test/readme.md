This folder contains scripts to facilitate testing of Prometheus rules and configuration of AlertManager.

| file | description |
|---|---|
| [download.sh](download.sh)| downloads prometheus, alertmanager and goreman for Mac OS. |
| [prometheus.yml](prometheus.yml) | Prometheus server configuration file. |
| [etcd_rules.yml](etcd_rules.yml) | Prometheus rules for etcd service up/down status changes. |
| [alertmanager.yml](alertmanager.yml) | AlertManager configuration file. |
| [Procfile](Procfile) | [Goreman](https://github.com/mattn/goreman) uses the Heroku Procfile to launch all required processes. |

Once *etcd, prometheus and alertmanager* are downloaded run the following command to start them up:
```bash
$ goreman start
```