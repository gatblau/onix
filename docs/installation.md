# Installing Onix CMDB

This section discusses a few different ways of deploying Onix CMDB services.

<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Installing using Docker Compose](#installing-using-docker-compose)
- [Installing using OpenShift](#installing-using-openshift)
- [Installing outside of containers](#installing-outside-of-containers)

<a name="installing-using-docker-compose"></a>
### Installing using Docker Compose [(up)](#toc)

[Docker Compose](https://docs.docker.com/compose/overview/) is a tool for defining and running multi-container Docker applications. 

In order to configure Onix's services, a docker compose YAML configuration file like [this](../install/container/docker-compose.yml) can be used.

In order to install Onix using Compose, a docker host is required. 
The current installation is for demo purposes only and runs two containers on the same host: the Onix database using a named data volume for persistent storage and the Keycloak server.
The Onix web service providing an API currently runs from the java runtime on the host.

**NOTE**: There is an issue when running Onix service in a container which prevents it from successfully authenticating with Keycloak. Further investigation is required.

To install Onix using this method follow the steps below:
- Ensure you have [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/overview/) installed on the Docker host.
- Execute the commands below:
```bash
$ mkdir onix && cd onix
$ wget https://raw.githubusercontent.com/gatblau/onix/master/install/container/docker-compose.yml
$ wget https://raw.githubusercontent.com/gatblau/onix/master/install/container/up.sh
$ sh up.sh {{onix-tag}}
$ java -jar -DAUTH_ENABLED=true onix-svc.jar
```

where the **{{onix-tag}}** is a combination of the git abbreviated last commit hash, and the time of the build:
 
 **[abbreviated last commit hash].[ddmmYY-hhMMss]** (e.g. 692e36a.220418-190211)
 
To find the appropriate tag take a look at the images available from Docker Hub:
- [Onix-Svc](https://hub.docker.com/r/gatoazul/onix-svc/) image repository.
- [Onix-Db](https://hub.docker.com/r/gatoazul/onix-db/) image repository.

**NOTE**: if calling docker-compose directly, mind the ONIXTAG environment variable must be exported as it is required by the [docker-compose.yml](../install/container/docker-compose.yml) file.

<a name="installing-using-openshift"></a>
### Installing using OpenShift [(up)](#toc)

Coming soon...

<a name="installing-using-openshift"></a>
### Installing outside of containers [(up)](#toc)

Installation outside of containers requires the preparation of the database and the java runtime environment for the web service.


#### Installing the database

The software requires PostgreSQL 9 or 10 server running.
Instructions to install PostgreSQL on all platforms can be found [here](https://www.postgresql.org/download/). 

If you are running on CentOS/RHEL distributions then a sample script for installing PostgreSQL client and server tools is provided [here](../install/vm/db/install_pgsql.sh).

Once the database server is up and running, it is necessary to configure the database schema for the CMDB.
[This script](../install/vm/db/prepare_onix_db.sh) automates the configuration of the database ready for the application to use.

#### Installing the web service

The web service is based on SpringBoot and provides a RESTful API running in an embedded web server in a fat Java Archive file.
The file can be built using [Apache Maven](https://maven.apache.org/) as follows:

```bash
$ git clone https://github.com/gatblau/onix.git
$ cd onix
$ mvn package
$ cd ..
$ cp ./onix/target/onix*.jar .
$ rm -rf onix
```

The fat jar created can be simply run using the **java** command as follows:

```bash
$ java -jar \
       -DHTTP_PORT=8080 \
       -DDB_USER=onix \
       -DDB_PWD=onix \
       -DDB_HOST=localhost \
       -DDB_PORT=5432 \
       -DDB_NAME=onix \
       -DAUTH_SERVER_URL='http://localhost:8081/auth' \
       -DAUTH_REALM=onix \
       -DAUTH_RESX=onix-cmdb \
       -DAUTH_ENABLED=true \
       onix-version.jar 
```
Where the following configuration variables are available to configure the jar file:

| Variable  | Description  | Default  |
|---|---|---|
| **HTTP_PORT** | the port number web service is listening on. | 8080  |
| **DB_USER**  | the user the service uses to connect to the database.  | onix  |
| **DB_PWD**  | the database user password.  | onix  |
| **DB_HOST**  | the database server hostname.  | localhost  |
| **DB_PORT**  | the posrt number the database server is listening on.  | 5432  |
| **DB_NAME**  | the name of the cmdb database the service attempts to connect.  | onix  |
| **AUTH_SERVER_URL** | The location of Keycloak auth endpoint. | http://localhost:8081/auth |
| **AUTH_REALM** | The application security realm. | onix |
| **AUTH_RESX** | The application to be authenticated (client_id). | onix-cmdb |
| **AUTH_ENABLED** | Whether IDAM is enabled. | false | 

