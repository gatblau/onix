# Binary deploy

The first step is to clone and build the Sentinel binary. Sentinel uses [modules](https://blog.golang.org/using-go-modules) to simplify dependency management, so building the binary will automatically download the required [dependencies](../go.mod).

Open the terminal and navigate to a folder of your preference where you want to download the source code, then type:

```bash
# downloads the source code and navigates to the created folder
git clone http://sentinel.gatblau.org && cd sentinel.gatblau.org

# make sure the go environment is set up correctly
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin

# build the binary
make build

# you should now see a binary file called sentinel in your current directory
ls sentinel
sentinel

# test you have kubectl and minikube
which kubectl && which minikube
/usr/local/bin/kubectl
/usr/local/bin/minikube

# start minikube using version 1.11.x - used by Sentinel
minikube start --kubernetes-version v1.11.10

# to test minikube is running
minikube status
host: Running
kubelet: Running
apiserver: Running
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.100

# check the content of the sentinel configuration
less config.toml

# ENSURE the following is set in the config.toml file:

# enables full trace to be logged
LogLevel = "Trace"

# sets the publisher to Logger
Publishers.Publisher = "logger"

# tell it to log to the file system
Publishers.Logger.OutputTo = "file"

# tell it where to log it
Publishers.Logger.LogFolder = "logs"

# now you can run sentinel
./sentinel

# sentinel should start logging and creating files in the logs folder
INFO[0000] Loading configuration.
INFO[0000] TRACE has been set as the logger level.       platform=kube-01
INFO[0000] LOGGER publisher has been registered.         platform=kube-01
...
...
```

Keep the above terminal session running and open a new terminal to do a deployment in minikube.
To test Sentinel we will deploy the **hello-minikube** example in the Kubernetes [quickstart documentation](https://kubernetes.io/docs/setup/minikube/#quickstart):

```bash
# deploy hello-minikube
kubectl run hello-minikube --image=k8s.gcr.io/echoserver:1.10 --port=8080

# switching to the folder where Sentinel is logging files
cd ~/sentinel.gatblau.org/logs

# you should be able to list the log files created
ls
1559984076799915000_deployment_CREATE_hello-minikube.json		1559984076863889000_deployment_UPDATE_hello-minikube.json
1559984076812404000_deployment_UPDATE_hello-minikube.json		1559984076871230000_pod_UPDATE_hello-minikube-59ddd8676b-vrdrs.json
1559984076822607000_deployment_UPDATE_hello-minikube.json		1559984078809792000_pod_UPDATE_hello-minikube-59ddd8676b-vrdrs.json
1559984076830423000_pod_CREATE_hello-minikube-59ddd8676b-vrdrs.json	1559984078822553000_deployment_UPDATE_hello-minikube.json
1559984076839006000_pod_UPDATE_hello-minikube-59ddd8676b-vrdrs.json

less 1559984076830423000_pod_CREATE_hello-minikube-59ddd8676b-vrdrs.json
```

```json
{
  "Change": {
    "name": "hello-minikube-59ddd8676b-vrdrs",
    "type": "CREATE",
    "namespace": "default",
    "kind": "pod",
    "time": "2019-06-08T08:54:36.830399Z",
    "host": "kube-01"
  },
  "Meta": {
    "name": "hello-minikube-59ddd8676b-vrdrs",
    "generateName": "hello-minikube-59ddd8676b-",
    "namespace": "default",
    "selfLink": "/api/v1/namespaces/default/pods/hello-minikube-59ddd8676b-vrdrs",
    "uid": "0b06686a-89cb-11e9-9aa9-08002796d2ab",
    "resourceVersion": "169030",
    "creationTimestamp": "2019-06-08T08:54:36Z",
    "labels": {
      "pod-template-hash": "1588842326",
      "run": "hello-minikube"
    },
    "ownerReferences": [
      {
        "apiVersion": "apps/v1",
        "kind": "ReplicaSet",
        "name": "hello-minikube-59ddd8676b",
        "uid": "0b03cc51-89cb-11e9-9aa9-08002796d2ab",
        "controller": true,
        "blockOwnerDeletion": true
      }
    ]
  }
}
```

```bash
# can delete deployment now
kubectl delete deployment hello-minikube

# more log files will be produced in the ./logs folder
```

Switch back to the other terminal and you should see that Sentinel trace output has been updated with details of the depployment (if LogLevel = "Trace" in the config.toml file).

```bash
# stops the Sentinel by pressing ctrl+c or cmd+c in macOS
```
