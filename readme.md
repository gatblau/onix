# Onix 

Onix is a configuration management database (CMDB) exposed via a RESTful API.

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

To run the service simply do:
```bash
$ java -jar -DDB_USER=onix -DDB_PWD=onix -DDB_HOST=localhost -DDB_PORT=5432 -DDB_NAME=onix onix-1.0-SNAPSHOT.jar 
```
where:
- DB_USER: database username
- DB_PWD: database user password
- DB_HOST: database server hostname
- DB_PORT: database server port
- DB_NAME: database name

#### On Container
- For an example of how to install the web service on a Docker container see [here](install/container/svc/build.sh).

## Web API Documentation

Onix uses Swagger to document its web API.

To see the Swagger UI go to http://localhost:8080/swagger-ui.html.

To see the API documentation in JSON format go to http://localhost:8080/v2/api-docs.

## Testing the service

#### Querying the service for version information

```bash
# replace the password with the password for the user
$ curl user:password@localhost:8080/
```

#### Creating a configuration item
```json
# create a payload.json file with the following content
# NOTE: the 'meta' value can be any json object, in this case is empty {}
# Use the 'meta' value to describe any specific property of your configuration item.
{
  "name": "Test Item",
  "description": "This is a CMDB item for testing purposes.",
  "itemTypeId": "2",
  "meta": "{ }",
  "tag": "Test",
  "deployed": false,
  "dimensions": [
    { "WBS" : "012csl" },
    { "COMPANY" : "ACME" }
  ]
}
```

```bash
# execute the PUT operation on the item URI passing a natural key and the payload.json file
$ curl -X PUT "user:password@localhost:8080/item/my_item_key" -F "payload.json"
```

#### Retrieving the configuration item using the natural key

```bash
# execute the GET operation on the item URI passing its item natural key
$ curl "user:password@localhost:8080/item/my_item_key" 
```

