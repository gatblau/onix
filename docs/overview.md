# Overview 

This section provides use cases for the CMDB and a technology overview.

<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Use Case Overview](#basic-use-cases)
- [Technology Overview](#technology-overview)

## Use Case Overview

The diagram below shows high level use cases for Onix based on different personas: 

![Onix Data Model](./pics/onix_arc.png "Onix Architecture") 

### Automation developer [1] [(up)](#toc)
 
Is concerned with **recording and testing configuration data changes** in a seamless way.

Developers write automation scripts and use connectors to record configuration data changes in the Onix CMDB via its Web API.
CMDB change records are tested as part of the automation development life cycle.

### Security officer [2] [(up)](#toc)
 
Is concerned with **securing access to configuration data** whilst leveraging organisation wide Identity and Access Management (IDAM) solutions.

To achieve this Onix provides [OIDC](https://openid.net/connect/) support.

### Operations team [3] [(up)](#toc)

Is concerned with **recording configuration data changes seamlessly** whilst executing automation scripts. 

As the development phase embedded and tested configuration data recording in the automation scripts, there is nothing for the Operations team to do other than execute the automation and observe changes made to the CMDB.

### End user [4] [(up)](#toc)

Is concerned with **requesting catalogue items** to be promptly deployed.

Any IT Self Service portal can provide the means to issue configuration change requests to an Automation Web API.
Configuration data changes are automatically recorded by the automation scripts.

### Data Analyst [5] [(up)](#toc)

Is concerned with **querying accurate configuration information at all times** for a variety of purposes.

As configuration data is updated when the automation scripts are executed, and has been tested as part of the development lifecycle, the information should be accurate and ready for reporting at all times.

### Development Project Members [6] [(up)](#toc)

Are concerned with **having visibility of applications and services** deployed on the infrastructure.


## Technology Overview

### Web API Service [(up)](#toc)

The [Web API](./wapi.md) is built using [Spring Boot](https://spring.io/projects/spring-boot).
It is a stateless RESTful style web service provided as a docker container image from Docker Hub.

The [WebAPI](./wapi.md) uses [Swagger](https://swagger.io/) for online documentation and testing of the service endpoints.

Access to the database is done using direct JDBC connections to ensure speed. All data logic is implemented as PostgreSQL functions.

### Database Service [(up)](#toc)

The database is implemented using [PostgreSQL](https://www.postgresql.org/).

### User Interface Service [(up)](#toc)

The user interface service is developed as a client of the Web API using [Vue.JS](https://vuejs.org/), [Nuxt.JS](https://nuxtjs.org/) and [Node.JS](https://nodejs.org/)

### Access Management [(up)](#toc)

The [Web API](./wapi.md) can be secured using Basic Authentication or OIDC.

### Deployment Configuration [(up)](#toc)

The solution requires an instance of PostgreSQL that can host the Onix Database.

The Web API and User Interface services should be deployed from Docker Images either on VM infrastructure (e.g. using Docker-Compose) or in a container platform such as OpenShift or Kubernetes using a deployment configuration.
