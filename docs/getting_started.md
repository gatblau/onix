# Getting started <img src="./pics/ox.png" width="200" height="200" align="right">
Getting started with Onix is easiy and required basically only two steps:

## Step 1: Deploy Container Images
The easiest way to get started is to deploy Onix from containers. Ready to use container images are provided on [docker hub](https://hub.docker.com/r/southwinds/) for each of the components of Onix (e.g. database, web api, user interface). For each component there are two repositories. One containing the stable GA release versions, the other containing latest snapshot releases.

There are different options to get started easily:
<a name="installing-using-docker"></a>
### Option A: Deploy Onix via Docker Compose
To do that you will need the following dependencies in your system:
- [Docker](https://docs.docker.com/compose/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Access to [Docker Hub Onix image repositories](https://hub.docker.com/u/gatblau)

Run the following commands from a terminal window:

```bash
# first, create a directory to place the docker-compose file
$ mkdir onix && cd onix

# second, download the docker compose file and place it into the new directory
$ wget https://raw.githubusercontent.com/gatblau/onix/v1/install/container/docker-compose.yml

# finally, launch docker compose in detached mode
$ docker-compose up -d
```

<a name="installing-using-openshift"></a>
### Option B: using OpenShift
Deploy using OpenShift is straightforward. 

Create an OpenShift project and deploy using one of the provided [templates here](install/openshift/readme.md).


<a name="installing-using-helm"></a>
### Option C: using Kubernetes and Helm Charts
To install Onix in plain Kubernetes, a [Helm Chart](https://helm.sh/docs/developing_charts/) will be provided soon.


## Step 2:  Initial contact using the Web API service
The easiest way to trying the service is via the Swagger API as describe below.

1. Open the Swagger UI (using the endpoint   `<YourHost>:<YourPort>/swagger-ui.html)` and navigate to the createOrUpdate PUT operation. E.g using docker on localhost, the link would be http://localhost:8080/swagger-ui.html#/web-api/createOrUpdateItemTreeUsingPUT. For OpenShift, replace with `<YourHost>:<YourPort>` with the route.
2. Paste the payload example found [here](../ansible/inventory/examples/inventory.json) in the payload box
3. Execute the request to create the inventory in the CMDB
4. You should see the response of the web service showing no errors

## Step 3: Check Image Configuration
The Container Images use env variables for their internal configuration.  Defaults should be okay for first experiments, but check that they fit to your environment.

## Onix Web API Image

The following variables are available to configure the Web API image:

| Variable  | Description  | Default  |
|---|---|---|
| **HTTP_PORT** | the port number web service is listening on. | 8080  |
| **DB_USER**  | the user the service uses to connect to the database.  | onix  |
| **DB_PWD**  | the database user password.  | onix  |
| **DB_HOST**  | the database server hostname.  | localhost  |
| **DB_PORT**  | the port number the database server is listening on.  | 5432  |
| **DB_NAME**  | the name of the cmdb database the service attempts to connect.  | onix  |
| **MGMT_ENDPOINT_METRICS_ENABLED** | enable metrics endpoints. | true |
| **DS_PREP_STMT_CACHE_SIZE** | number of prepared statements that the JDBC driver will cache per connection. | 250 |
| **DS_PREP_STMT_CACHE_SQL_LIMIT** | maximum length of a prepared SQL statement that the driver will cache. | 2048 |
| **DS_CACHE_PREP_STMTS** | enable the cache. | true |
| **DS_USE_SERVER_PREP_STMTS** | add support for server-side prepared statements. | true |
