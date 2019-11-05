<img src="./../pics/sentinel_small.png" align="right" height="200" width="200"/>

# Configuration

Sentinel can be configured by either updating variables in its [configuration file](../config.toml) or setting environment variables.

When environment variables are set, they override the values in the configuration file.

The configuration file and corresponding environment variables are described below.

## General Variables

| Config File Var | Environment Var | Description | Default |
|---|---|---|---|
| Platform | SL_PLATFORM | the identifier of the platform . | KUBE-01 |
| KubeConfig | SL_KUBECONFIG | the path to the kubernetes configuration file used by the Sentinel to connect to the kubernetes API. | ~/.kube/config |
| LoginLevel | SL_LOGINLEVEL | defines the login level used by the software. Possible values are: __Trace, Debug, Info, Warning, Error, Fatal and Panic__. | Info |
| Publishers.Publisher| SL_PUBLISHERS_PUBLISHER | defines which publisher to use (i.e. webhook, broker or logger). The logger publisher is there to write to standard output. | logger |

### _Logger Publisher Variables_

| File Var | Environment Var | Description | Default |
|---|---|---|---|
| Publishers.Logger.OutputTo | SL_PUBLISHERS_LOGGER_OUTPUTTO | whether to log to the standard output (stdout) or to the file system (file) | stdout |
| Publishers.Logger.LogFolder| SL_PUBLISHERS_LOGGER_LOGFOLDER | the path to the log folder, only required if Output = "file" | logs |

## _Webhook Publisher Variables_

The webhook publisher allows for the configuration of one or more target endpoints/web consumer applications. For this reason, whithin the [config.toml](../config.toml) file, each endpoint requires a separate [TOML](https://github.com/toml-lang/toml) table using double square brackets syntax: __[[Publishers.Webhook]]__ - the TOML definition of array of tables.

Variables within each TOML table, are mapped to environment variables using indices as follows:

| Environment Var | Table | Description |
|---|:-:|---|
|__SL_PUBLISHERS_WEBHOOK_0_URI__| 0 | URI of the first endpoint |
|__SL_PUBLISHERS_WEBHOOK_1_URI__| 1 | URI of the second endpoint |

__NOTE__: if the intention is to use environment variables to configure multiple webhooks, then as many _[[webhook]]_ tables need to be created in the config.toml file as web hooks are required to set in environment variables. This is because the binding of environment variables to configuration file variables is based on the number of _[[webhook]]_ tables in the [config.toml](../config.toml) file.

The following table shows all configuration variables available to the webhook publisher:

| File Var | Environment Var | Description | Default |
|---|---|---|---|
| Publishers.Webhook.URI | SL_PUBLISHERS_WEBHOOK_[X]_URI | the uri of the webhook | localhost:8080/sentinel |
| Publishers.Webhook.Authentication | SL_PUBLISHERS_WEBHOOK_[X]_AUTHENTICATION | authentication mode to use for posting events to the webhook endpoint (i.e. none, basic) | - |
| Publishers.Webhook.Username | SL_PUBLISHERS_WEBHOOK_[X]_USERNAME | the optional username for basic authentication | sentinel |
| Publishers.Webhook.Password | SL_PUBLISHERS_WEBHOOK_[X]_PASSWORD | the optional password for basic authentication | s3nt1nel |

## _Broker Publisher Variables_

| File Var | Environment Var | Description | Default |
|---|---|---|---|
| Publishers.Broker.Brokers | SL_PUBLISHERS_BROKER_BROKERS | the Kafka brokers to connect to, as a comma separated list | localhost:9092 |
| Publishers.Broker.Certificate | SL_PUBLISHERS_BROKER_CERTIFICATE | optional certificate file for client authentication | - |
| Publishers.Broker.Key | SL_PUBLISHERS_BROKER_KEY | optional key file for client authentication | - |
| Publishers.Broker.CA | SL_PUBLISHERS_BROKER_CA | optional certificate authority file for TLS client authentication | - |
| Publishers.Broker.Verify | SL_PUBLISHERS_BROKER_VERIFY | optional verify ssl certificates chain | false |

## _Observable Objects Flags_

The following flags can be used to switch on/off the objects Sentinel observes:

| File Var | Environment Var | Description | Default |
|---|---|---|---|
| Observe.Service | SL_OBSERVE_SERVICE | whether to observe create, update and delete service events | true |
| Observe.Pod | SL_OBSERVE_POD | whether to observe create, update and delete pod events | true |
| Observe.PersistentVolume | SL_OBSERVE_PERSISTENTVOLUME | whether to observe create, update and delete persistent volume events | true |
| Observe.PersistentVolumeClaim | SL_OBSERVE_PERSISTENTVOLUMECLAIM | whether to observe create, update and delete persistent volume claim events | true |
| Observe.Namespace | SL_OBSERVE_NAMESPACE | whether to observe create, update and delete namespace events | true |
| Observe.Deployment | SL_OBSERVE_DEPLOYMENT | whether to observe create, update and delete deployment events | false |
| Observe.ReplicationController | SL_OBSERVE_REPLICATIONCONTROLLER | whether to observe create, update and delete replication controller events | true |
| Observe.ReplicateSet | SL_OBSERVE_REPLICASET | whether to observe create, update and delete replica set events | false |
| Observe.DaemonSet | SL_OBSERVE_DAEMONSET | whether to observe create, update and delete daemon set events | false |
| Observe.Job | SL_OBSERVE_JOB | whether to observe create, update and delete job events | false |
| Observe.Secret | SL_OBSERVE_SECRET | whether to observe create, update and delete secret events | false |
| Observe.ConfigMap | SL_OBSERVE_CONFIGMAP | whether to observe create, update and delete config map events | false |
| Observe.Ingress | SL_OBSERVE_INGRESS | whether to observe create, update and delete ingress events | false |
| Observe.ServiceAccount | SL_OBSERVE_SERVICEACCOUNT | whether to observe create, update and delete service account events | false |
| Observe.ClusterRole | SL_OBSERVE_CLUSTERROLE | whether to observe create, update and delete cluster role events | false |
| Observe.ResourceQuota | SL_OBSERVE_RESOURCEQUOTA | whether to observe create, update and delete resource quota events | true |
| Observe.NetworkPolicy | SL_OBSERVE_NETWORKPOLICY | whether to observe create, update and delete network policy events | false |


[*] _The Sentinel icon was made by [Freepik](https://www.freepik.com) from [Flaticon](https://www.flaticon.com) and is licensed by [Creative Commons BY 3.0](http://creativecommons.org/licenses/by/3.0)_