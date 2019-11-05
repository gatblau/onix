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
| Id | OXKU_ID | the oxkube service instance Id for logging identification purposes. | OxKube-01 |
| LogLevel | OXKU_LOGLEVEL | verbosity of logging (Trace, Debug, Warning, Info, Error, Fatal, Panic) | Trace |
| Metrics | OXKU_METRICS | enables metrics | true |
| Consumers.Consumer | OXKU_CONSUMERS_CONSUMER | - | - |

### Onix Configuration

| Cfg Vars | Env Vars | Description | Default |
|---|---|---|---|
| Onix.URL | OXKU_URL | The URL of the Onix WAPI service to connect to. | http://localhost:8080 |
| Onix.AuthMode | OXKU_AUTHMODE | the athentication type to use when connecting to the Onix WAPI service (none, basic or oidc) | basic |
| Onix.Username | OXKU_USERNAME | the user name if basic authentication mode is used. | admin |
| Onix.Password | OXKU_PASSWORD | the password if basic authentication mode is used. | 0n1x |
| Onix.ClienId | OXKU_CLIENTID | The Client Id if OpenId Connect authentication mode is used. | - |
| Onix.ClienSecret | OXKU_CLIENTSECRET | The client secret if OpenId Connect authentication mode is used. | - |
| Onix.TokenURI | OXKU_TOKENURI | The Token service URI if OpenId Connect authentication mode is used. | - |

### Web Hook Consumer Configuration

| Cfg Vars | Env Vars | Description | Default |
|---|---|---|---|
| Consumers.Consumer.Webhook.Port | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PORT | - | - |
| Consumers.Consumer.Webhook.Path | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PATH | - | - |
| Consumers.Consumer.Webhook.AuthMode | OXKU_CONSUMERS_CONSUMER_WEBHOOK_AUTHMODE | - | - |
| Consumers.Consumer.Webhook.Username | OXKU_CONSUMERS_CONSUMER_WEBHOOK_USERNAME | - | - |
| Consumers.Consumer.Webhook.Password | OXKU_CONSUMERS_CONSUMER_WEBHOOK_PASSWORD | - | - |
| Consumers.Consumer.Webhook.InsecureSkipVerify | OXKU_CONSUMERS_CONSUMER_WEBHOOK_INSECURESKIPVERIFY | - | - |
