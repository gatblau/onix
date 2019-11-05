<img src="../../pics/sentinel_small.png" align="right" height="200" width="200"/>

# Kubernetes deployment using helm

In order to deploy in Kubernetes, do the following:

1. Install the [helm client](https://helm.sh/docs/using_helm/#installing-helm).

2. Tiller is the server-side component for helm. Tiller needs to be present in the kubernetes cluster to deploy applications using helm charts. If Tiller is not already installed, install it by running the [install_tiller.sh](install_tiller.sh) script as a cluster-admin.

3. Sentinel requires cluster level permission to watch changes of Kubernetes resources. The [setup_rbac.sh](setup_rbac.sh) script should be executed to configure the privileges to run.  

4. Finally, Sentinel can be installed using the provided helm chart, runningthe following command:

```bash
helm install sentinel \
  --name sentinel-release \
  --namespace sentinel \
  --values sentinel/values.yaml
```

___Note___: the [values.yaml](sentinel/values.yaml) file contains the default configuration for the chart.
It needs to be updated before running the chart with the correct publisher connection settings. This is so that Sentinel can connect to a Webhook or a message broker to publish change information.

[[back to index](../readme.md)]

[*] _The Sentinel icon was made by [Freepik](https://www.freepik.com) from [Flaticon](https://www.flaticon.com) and is licensed by [Creative Commons BY 3.0](http://creativecommons.org/licenses/by/3.0)_