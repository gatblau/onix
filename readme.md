# Onix 

## Installation Notes

### Installing the database

The software requires PostgreSQL server running.

#### On VM
- For an example of how to install the database on RHEL/CentOS VM see [here](install/vm/db/install_pgsql.sh).
- For an example of how to configure the database schema for the CMDB see [here](install/vm/db/configure_pgsql.sh).

#### On Container
- For an example of how to create a Docker container for the database see [here](install/container/db/build.sh).

### Installing the web service

The web service is a SpringBoot Restful service running an embedded web server in a uber jar.

#### On VM
- For an example of how to install the web service on RHEL/CentOS VM see [here](install/vm/svc/build.sh).

#### On Container
- For an example of how to install the web service on RHEL/CentOS VM see [here](install/container/svc/build.sh).

