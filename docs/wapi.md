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
| Method | PUT | 
| Path | /item/**{item_key}**/|
| Response Content Type | application/json |
 

### Sample Payload

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item is not part of the payload but specified in the URI.

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

### Payload fields

The following table describes the fields in the payload and provides some examples of use:

| Field  | Description  | Example |
|---|---|---|
| name | The name associated with the configuration item. | "OCP Demo Master 1" |
| description | A description of the configuration item. | "OpenShift Demo platform first master in zone A."
| itemTypeId | The unique Id of the configuration item type. It must exist as a valid item type. | "2" |
| meta | Stores any well-formed json object. This is the primary mechanism to store configuration item information. | { "host":"OCP-DEMO-M-01", "region":"Ireland", "provider":"AWS" } |
| tag | Used for annotating the item for searching. For example, a search can be done by items having the EUROPE tag.| "TEST RELEASE-B EUROPE" |
| deployed | Indicates if this item has been deployed or is waiting to be deployed. | true/false |
| dimensions | A JSON array of key and value pairs used for reporting. | As per sample payload above. |

### Example

The following example shows how to execute a PUT request to the service using [cURL](https://curl.haxx.se/):

```bash
# execute the PUT operation on the item URI passing a natural key (e.g. KEYDEMOM001) and a payload via json file with contents as per sample above.
$ curl -X PUT -H 'ContentType: application/json' -d '@item_payload.json' 'user:password@localhost:8080/item/KEYDEMOM001' 
```

## Retrieving the configuration item by key [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /item/**{item_key}**/|
| Response Content Type | application/json or application/x-yaml |

```bash
# execute the GET operation on the item URI passing its natural key
$ curl 'user:password@localhost:8080/item/KEYDEMOM001' 
```

## Linking two items [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | PUT | 
| Path | /link/**{link_key}**/|
| Response Content Type | application/json |

### Sample Payload

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item is not part of the payload but specified in the URI.

```json
{
  "meta": "{ }",
  "description": "This is a CMDB item for testing purposes.",
  "role": "connect",
  "start_item_key": "ITEM_ONE_KEY",
  "end_item_key": "ITEM_TWO_KEY",
  "tag": "Test"
}
```

### Payload fields

The following table describes the fields in the payload and provides some examples of use:

| Field  | Description  | Example |
|---|---|---|
| description | A description of the link. | "Link A to B."
| meta | Stores any well-formed json object. This is the primary mechanism to store link configuration information. | { "key1":"value1", "key2":"value2", "key3":"value3" } |
| tag | Used for annotating the link for searching. For example, a search can be done by links having the TEST tag.| "TEST RELEASE-B EUROPE" |
| role | Indicates the role of the link. | "is installed on", "connects to" |
| start_item_key | The key of the item from which the link starts. | A configuration item key. |
| end_item_key | The key of the item to which the link ends. | A configuration item key. |

### Example

The following example shows how to execute a PUT request to the service using [cURL](https://curl.haxx.se/):

```bash
# execute the PUT operation on the item URI passing the link natural key and the payload.json file
$ curl -X PUT -H 'ContentType: application/json' -d '@link_payload.json' 'user:password@localhost:8080/link/my_link_key/' 
```

## Retrieving a link [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /link/**{link_key}**/|
| Response Content Type | application/json or application/x-yaml |

### Example

```bash
# execute the GET operation on the item URI passing the link natural key 
$ curl 'user:password@localhost:8080/link/my_link_key/' 
```
