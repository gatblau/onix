# Launching & Testing Buildman

To test buildman it is assumed minikube is running locally.
Alternatively, a server side K8S flavour can aslo be used.

**NOTE** Artisan CLI is required to run the commands below.

## Installing Minikube

If you dont have a testing platform you can use minikube.

To install it run:

```bash
# in windows
art run install-mini-win

# in linux
art run install-mini-linux

# in MacOS
art run install-mini-mac

# then start it
art run start-mini
```

## Installing Tekton

Buildman uses Tekton, to install it run:

```bash
# in windows
art run install-tekton

# once the tailed logs show the 2 containers running then CTRL+C to stop logs
```

## Setting up RBAC

In order for Buildman to work it needs access to tekton.dev api group.

To set up the necessary access run:

```bash
art run setup-rbac
```

The above should create a service account and a bound namespaced role.

**NOTE**: By default they are created in the default namespace.

If you have your own namespace where Buildman is to run, then update the NAMESPACE variable in [build.yaml](build.yaml)

## Deploy Buildman

To deploy buildman, run:

```bash
art run buildman-deploy
```

This will create a config map with a single test Buildman policy in [buildman.yaml](buildman.yaml).

If you want to try your specific policies then update the content of the config map.

**NOTE**: By default buildman is deployed in the *default* namespace.

If you have your own namespace where Buildman is to run, then update the NAMESPACE variable in [build.yaml](build.yaml)

## Outcome

Buildman detects changes in a dummy image and dummy base in quay.io and triggers a new pipleine run every minute.

The pipeline run fails as no pipeline is installed but shows buildman launching pipeline runs in K8S.

## Image build policy

Buildman uses policies in the configmap to decide whether to build an image.

Each policy has the following attributes:

|name| description |
|---|---|
|`name`| the policy name that should match the name of the name of the pipelinerun to be created - e.g. **dummy-app** -> **dummy-app**-image-pr-xxxx |
|`description`|a description of the policy - what it does |
|`app`| the name of the application image to monitor |
| `app-user` | the username used to connect to the app image registry - if empty then assume no authentication |
| `app-pwd` | the password used to connect to the app image registry - if empty then assume no authentication |
| `app-base-created-label`| the name of the label in the application image which contains the build date of the date image - e.g. '*build-date*' |
|`base`||
|`base-user`| the username used to connect to the base image registry - if empty then assume no authentication |
|`base-pwd`| the password used to connect to the base image registry - if empty then assume no authentication |
| `namespace`  | the name of the kubernetes namespace where pipelineruns are created |
| `pollBase` | true if Buildman should poll the base image for changes |
