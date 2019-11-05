# Installing on OpenShift

## Configuring Security

Sentinel requires *watch* and *list* privileges on the resources is watching for changes across the Kubernetes cluster. These privileges are described in a new [cluster role](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#role-and-clusterrole) that can be found [here](cluster_role.yaml).

Before running the template, the above cluster role needs to be created with the privileges required by Sentinel to run.

The role then needs to be bound to the sentinel service account.

To create the cluster role, service account, role binding and namespace, log in as cluster admin and execute the following:

```bash
# setup the necessary security for Sentinel access to Kubernetes resources
$ make oc-setup
```

## Installing using the web console catalogue

Import the Sentinel template as follows:

```bash
# import the template in the catalogue
$ make oc-import-template
```

Once the template is imported in OpenShift, it shows in the catalogue and can be run using the web console.

### **PLEASE NOTE!**

The provided template deploys Sentinel using the **stdout** option of the logger publisher by default.

The logger publisher is only used for testing Sentinel without really publishing events but writing them out to the logs - or alternatively to the file system if Sentinel is running outside of Kubernetes.

To publish events to downstream systems *a different logger must be used*. This can be done by modifying the required environment variables for the selected publisher when running the template from the web console.

Below is an example of the minimum variables required by the **webhook** publisher:

| variable | example |
|---|---|
| **PUBLISHER** | webhook |
| **WEBHOOK_URI** | https://change-to-the-onix-kube-uri/webhook |
| **WEBHOOK_USER** | sentinel |
| **WEBHOOK_PWD** | kfljetvjn |


## Installing using the command line

```bash
# you are log as system admin!
$ oc login -u system:admin

# create a new project, role, account and bindings
$ sh ./scripts/openshift/setup.sh

# deploy the app from the file system passing in 
# the publisher connection information
# update the publisher connection information to suit your needs
$ oc new-app ./scripts/openshift/sentinel.yml -p PUBLISHER=webhook -p WEBHOOK_URI=https://change-to-the-onix-kube-uri/webhook -p WEBHOOK_USER=change_me_user -p WEBHOOK_PWD=change_me_pwd
```

## Cleanup operations

To delete the template:

```bash
# using make
$ make oc-delete-template

# using oc
$ oc delete template sentinel -n openshift
```

To remove all Sentinel resources run the command below which ensures that in addition to the Sentinel namesapce, the **clusterrolebinding** and **clusterrole** are also removed:

```bash
$ make oc-cleanup
```
