<img src="./pics/ox.png" width="200" height="200" align="right">

# Database Deployment

The Onix database is based on the [PostgreSQL](https://www.postgresql.org/) open source object-relational database system.

In order to simplify the deployment and the upgrades of the database schemas and functions, Onix has a built-in algorithm that deals with these database operations.

## Initial Deployment

When Onix is deployed for the first time, it is configured to connect to a clean instance of a PostgreSQL server. Typically, and as Onix is designed to run in a container platform like [OpenShift](https://www.openshift.com/) or [Kubernetes](https://kubernetes.io/), it triggers the deployment of the database when the [readiness probe web api endpoint](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/#define-readiness-probes) (e.g. http://localhost:8080/ready URL) is called.

The readyness algorithm perform the following checks:

1. Checks if the Onix database exists, if so go and check for upgrade requirements (see next section below)
2. If not, creates the Onix database and user
3. Gets the database schema version from the application manifest
4. Fetches the database deployment scripts from either a local (embedded in the container image) or remote (github.com) source.
5. Connects to the database and runs the scripts
6. Updates the version control table
7. Check for upgrade requirements

## Upgrades

The auto upgrade logic is as follows:

1. Compares the current database schema version against the one required by the current application, then
2. If the current database schema version is the same as the version required by the application then go to ready state
3. If the current database schema version is greater than the version required by the application then raises an exception - the container will not run
4. If the current database schema version is smaller than the version required by the application then triggers the update procedure if auto updates is enabled 

### Upgrade Procedure

Currently, the upgrade procedure logic is not implemented (i.e. raises Unsupported Operation Exception). However, the logic when implemented, will be as follows:

1. Drop all functions for the current schema version
2. Execute schema changes and data migrations
3. Re-create all functions for the new schema version

------

## Database Scripts Folder Tree

The following structure is required for the deploy/upgrade algorithm to select the relevant upgrade scripts:

```ANSI
./ { root of the structure }
+-- app/ {app manifests folder}
|   +-- 0.0.1.json {app manifest for version 0.0.1}
|   +-- 0.0.2.json {app manifest for version 0.0.2}
|   +-- 0.0.3.json {app manifest for version 0.0.2}
+-- db/ {database scripts folder}
|   +-- install/ { installation scripts folder }
|       +-- 1/ { schema version 1 forlder }
|           +-- db.json { db manifest for schema version 1 }
|           +-- script_1.sql { sql script }
|           +-- script_2.sql { sql script }
|           +-- script_3.sql { sql script }
|       +-- 2/ { schema version 2 forlder }
|           +-- upgrade/ { upgrades from previous version (1) script folder }
|               +-- script.sql
|           +-- db.json { db manifest for schema version 2 }
|           +-- script_1.sql { sql script }
|           +-- script_2.sql { sql script }
|           +-- script_3.sql { sql script }
|       +-- 3/ { schema version 3 forlder }
|           +-- upgrade/ { upgrades from previous version (2) script folder }
|               +-- script.sql
|           +-- db.json { db manifest for schema version 3 }
|           +-- script_1.sql { sql script }
|           +-- script_2.sql { sql script }
|           +-- script_3.sql { sql script }

```

**NOTE**: *the structure above is embedded as resources in the application jar file in the docker image. It can also be deployed on a github repository for remote reading although it is not the preferred storage approach and might be deprecated in future version*.

----------

## Environment Variables

The following variables control the behaviour of the database operational algorithm:

| Name | Description | Default |
|---|---|---|
| **DB_SCRIPTS_REMOTE** | Whether to use remote scripts (deployed on remote git). The default behaviour is to use local scripts (embedded in the image) | false |
| **DB_SCRIPTS_URL** | The URL of the folder structure root in the remote git repository. | "https://raw.githubusercontent.com/ gatblau/onix/ <app_commit>/src/main/resources". Note that **<app_commit>** tag in the URL is automatically replaced by the commit hash of the application release in git. |
| **DB_AUTO_DEPLOY** | Whether the readyness probe will attempt to auto deploy the database if it does not exist. | true |
| **DB_AUTO_UPGRADE** | Whether the readyness probe will attempt to auto upgrade the database if its version is lower than the one required by the current application. | false |
