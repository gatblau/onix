# Getting started

## Installing the required services

The easiest way to get started is to deploy the Web API and Database services from containers.

To do that you will need the following dependencies in your system:

- [Docker](https://docs.docker.com/compose/install/)
- [Docker Compose](https://docs.docker.com/compose/install/)
- Access to [Docker Hub Onix image repositories](https://hub.docker.com/u/southwinds)

Run the following commands from a terminal window:

```bash
# first, create a directory to place the docker-compose file
$ mkdir onix && cd onix

# second, download the docker compose file and place it into the new directory
$ wget https://github.com/gatblau/onix/blob/v1/install/container/docker-compose.yml

# finally, launch docker compose in detached mode
$ docker-compose up -d
```

<a name="installing-using-openshift"></a>
### Installing using OpenShift [(up)](#toc)

To install Onix in OpenShift or Kubernetes, a [Helm Chart](https://helm.sh/docs/developing_charts/) will be provided soon.


### Image Configuration

#### Onix Web API Image

There are two repositories for the Web API image:
- [Release Repository](https://hub.docker.com/r/southwinds/onixwapi): used for GA releases.
- [Snapshot Repository](https://hub.docker.com/r/southwinds/onixwapi-snapshot): used for snapshot releases.

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

## Trying the Web API service

The easiest way to trying the service is via the Swagger API as describe below.

1. Open the [Swagger UI](http://localhost:8080/swagger-ui.html#/web-api/createOrUpdateItemTreeUsingPUT)
2. Paste the payload [here](./../connectors/ansible/inventory/examples/inventory.json) in the payload box
3. Execute the request to create the inventory in the CMDB
4. You should see the response of the web service showing no errors


