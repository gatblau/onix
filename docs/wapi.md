# Web API Documentation 

This section explains how to use Onix Web API. 


<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Using Swagger](#using-swagger)
- [Getting Service Information](#getting-service-information)


<a name="using-swagger"></a>
## Using Swagger [(up)](#toc)

Onix uses [Swagger](https://swagger.io/) to document its web API. 

### WAPI User Interface

When Onix is up and running, the Swagger User Interface can be reached at the following URI:

**http://onix_host_name:onix_port_number/swagger-ui.html**
 
### JSON WAPI Docs

Similarly, a JSON representation of the Web API documentation can be obtained from the following URI: 
 
**http://onix_host_name:onix_port_number/v2/api-docs**


<a name="getting-service-information"></a>
## Getting Service Information [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | / |
| Response Content Type | application/json |
 

### Usage example:

```bash
# replace the password with the password for the user
$ curl user:password@localhost:8080/
```

## Creating a configuration item [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /item/**{item_key}**/|
| Response Content Type | application/json |
 

### Sample Payload:

**NOTE**: the 'meta' value can be any json object, in this case is empty {}.
Use the 'meta' property value to describe any specific property of your configuration item.

```json
{
  "name": "This is the name associated with the configuration item.",
  "description": "This a description of the configuration item.",
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

### Usage example:

**NOTE**: 

```bash
# execute the PUT operation on the item URI passing a natural key and the payload.json file
$ curl -X PUT "user:password@localhost:8080/item/my_item_key" -f "item_payload.json"
```

#### Retrieving the configuration item using the natural key

```bash
# execute the GET operation on the item URI passing its natural key
$ curl "user:password@localhost:8080/item/my_item_key" 
```

#### Creating a link between two items

```json
# Create a link_payload.json file with the following content
# NOTE: the 'meta' value can be any json object, in this case is empty {}
# Use the 'meta' value to describe any specific property of your configuration item.
{
  "meta": "{ }",
  "description": "This is a CMDB item for testing purposes.",
  "role": "connect",
  "start_item_key": "ITEM_ONE_KEY",
  "end_item_key": "ITEM_TWO_KEY",
  "tag": "Test"
}
```

```bash
# execute the PUT operation on the item URI passing the link natural key and the payload.json file
$ curl -X PUT "user:password@localhost:8080/link/my_link_key/" -F "link_payload.json"
```

#### Retrieving a link between two items
```bash
# execute the GET operation on the item URI passing the link natural key 
$ curl "user:password@localhost:8080/link/my_link_key/" 
```
