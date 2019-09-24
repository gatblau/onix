# OpenShift Templates <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

This folder contains various OpenShift templates for different deployment scenarios:

| Template | Description |
|---|---|
|[Onix with ephemeral storage](onix-ephemeral.yml) | Deploys Onix and PostgreSQL with [ephemeral storage](https://docs.openshift.com/online/architecture/additional_concepts/ephemeral-storage.html). |
|[Onix with persistent storage](onix-persistent.yml) | Deploys Onix and PostgreSQL with [persistent storage](https://docs.openshift.com/online/architecture/additional_concepts/storage.html). |
|[Onix with external database](onix-ext-db.yml)| Deploys Onix and connects it to an external PostgreSQL database. |
| [Onix Kubernetes Agent](oxkube.yml) | Deploys the Kubernetes Agent for Onix. |
| [Onix Web Console](oxwc.yml) | Deploys the Web Console only. |
| [Onix All Persistent](onix-all-persistent.yml) | Deploys all onix services (i.e. the database, web api, web console and kubernetes agent.) |

__NOTE__: Onix will automatically deploy the SQL schemas when the readyness probe is called upon deployment of the Onix WAPI container.

The *persistent template* uses an OpenShift volume claim as durable storage for the database component.  

The *ephemeral template* use an empty dir as storage for the database, ***which means all data is lost if the database container is restarted***. Use only if you don't have any persistent storage at hand.

The *external database template* uses an [external service](https://docs.openshift.com/online/dev_guide/integrating_external_services.html#mysql-define-service-using-fqdn) to connect to an instance of PostgreSQL that is external to OpenShift.

## Creating the application using the cli

The following steps use the **oc command line tool** to create a new empty project and deploy using the persistant template.

```bash
# first, create an empty project
$ oc new-project onix

# Deploy Onix using the persistent storage option
$  oc new-app https://raw.githubusercontent.com/gatblau/onix/master/docs/install/openshift/onix-all-persistent.yml
```

## Importing the templates

To [import the templates](import.sh) in OpenShift, run the following command:

```bash
# log in as admin
oc login -u system:admin

# import the templates so they show in the catalogue
sh import.sh
```

Once the templates have been imported, they should be visible in the Web Console under the Catalogue section in OpenShift.

## PostgreSQL cloud providers

The table below list a few cloud providers running PostgreSQL database services that could be used as a backing database for Onix:

| Provider | Service |
|---|---|
| AWS | [Amazon RDS](https://aws.amazon.com/rds/postgresql/) |
| Azure | [Azure Database for PostgreSQL](https://azure.microsoft.com/en-us/services/postgresql/) |
| Google Cloud | [Cloud SQL for PostgreSQL](https://cloud.google.com/sql/docs/postgres/) |
