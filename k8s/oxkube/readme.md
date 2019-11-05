<img src="../../docs/pics/ox.png" width="200" height="200" align="right">

# Onix API Extension for Kubernetes

Ox-Kube is an [Onix API Extension](../../docs/extensions.md) API extension for [Kubernetes](http://kubernetes.io).

 It consumes events sent by [Sentinel](http://sentinel.gatblau.org),
 either via a webhook or a message broker consumer, and updates the Onix database when the status of Kubernetes resources change.

![OxKube](docs/pics/ox_kube.png)

## Kubernetes Model

 The Onix model representing Kubernetes objects can be found [here](../../docs/models/k8s/readme.md).

## Configuration

 OxKube is configured via the [config.toml](src/config.toml) file.

 In addition, environment variables can be used to override the values in the configuration file as follows:

### General Configuration

| Cfg Vars | Env Vars | Description | Default |
|---|---|---|---|
| Id | OXKU_ID | - | - |
| LogLevel | OXKU_LOGLEVEL | - | - |
| Metrics | OXKU_METRICS | - | - |
| Consumers.Consumer | OXKU_CONSUMERS_CONSUMER | - | - |

### Onix Configuration

| Cfg Vars | Env Vars | Description | Default |
|---|---|---|---|
| Onix.URL | OXKU_ID | - | - |
| Onix.AuthMode | OXKU_AUTHMODE | - | - |
| Onix.Username | OXKU_USERNAME | - | - |
| Onix.Password | OXKU_PASSWORD | - | - |

### Web Hook Consumer Configuration

| Cfg Vars | Env Vars | Description | Default |
|---|---|---|---|
| Consumers.Consumer.Webhook.Port | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PORT | - | - |
| Consumers.Consumer.Webhook.Path | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PATH | - | - |
| Consumers.Consumer.Webhook.AuthMode | OXKU_CONSUMERS_CONSUMER_WEBHOOK_AUTHMODE | - | - |
| Consumers.Consumer.Webhook.Username | OXKU_CONSUMERS_CONSUMER_WEBHOOK_USERNAME | - | - |
| Consumers.Consumer.Webhook.Password | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PASSWORD | - | - |
| Consumers.Consumer.Webhook.InsecureSkipVerify | OXKU_CONSUMERS_CONSUMER_WEBHOOK_INSECURESKIPVERIFY | - | - |
