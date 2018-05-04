# Web API 

This section explains how to use Onix Web API. 


<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Using Swagger](#using-swagger)
- [Identity and Access Management](#identity-and-access-management)
- [Getting service information](#getting-service-information)
- [Creating a configuration item](#creating-a-configuration-item)
- [Retrieving a configuration item by key](#retrieving-a-configuration-item-by-key)
- [Linking two items](#linking-two-items)
- [Retrieving a link](#retrieving-a-link)

<a name="using-swagger"></a>
## Using Swagger [(up)](#toc)

Onix uses [Swagger](https://swagger.io/) to document its web API. 

### WAPI User Interface

When Onix is up and running, the Swagger User Interface can be reached at the following URI:

**http://onix_host_name:onix_port_number/swagger-ui.html**
 
### JSON WAPI Docs

Similarly, a JSON representation of the Web API documentation can be obtained from the following URI: 
 
**http://onix_host_name:onix_port_number/v2/api-docs**

<a name="identity-and-access-management"></a>
## Identity and Access Management [(up)](#toc)

If the Onix Service is configured with '**AUTH_ENABLED = true**', then a bearer token must be passed with every request via an authorisation header.
The following example shows how to obtain a token and pass it to the service:

```bash

# gets a token from the Auth server
# NOTE: replace the payload attributes depending on the configuration of the Auth server
TOKEN=`curl -d "client_id=onix-cmdb" -d "username=onix" -d "password=onix" -d "grant_type=password" "http://localhost:8081/auth/realms/onix/protocol/openid-connect/token"`

# constructs an authorization header
# NOTE: check the access_token substrings below is OK for your case
AUTH_HEADER='Authorization: bearer '${TOKEN:17:1135} 

# executes the request passing the bearer toke via the AUTH_HEADER
curl \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/itemtype/'
```

For more information on authentication see the [IDAM section](idam.md).

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
$ curl \
    -H '${AUTH_HEADER}' \
    'http:localhost:8080'
```
<a name="creating-a-configuration-item"></a>
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
  "meta": { 
    "a_custom_key" : "a_custom_value",
    "another_custom_key" : "another_custom_value",  
    "a_custom_complex_object" : {
      "a_custom_key" : "a_custom_value",
      "another_custom_key" : "another_custom_value"
    }
  },
  "tag": "Test",
  "deployed": false,
  "dimensions": { 
    "WBS" : "012csl", 
    "COMPANY" : "ACME" 
  }
}
```
**NOTE**: the *meta* field can contain *any* JSON object.

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
| dimensions | A set of of key and value pairs used for reporting. | As per sample payload above. |

### Example

The following example shows how to execute a PUT request to the service using [cURL](https://curl.haxx.se/):

```bash
# execute the PUT operation on the item URI passing a natural key (e.g. KEYDEMOM001) and a payload via json file with contents as per sample above.
$ curl \
    -X PUT \
    -H 'ContentType: application/json' \
    -H '${AUTH_HEADER}' \
    -d '@item_payload.json' \
    'http://localhost:8080/item/KEYDEMOM001' 
```

<a name="retrieving-a-configuration-item-by-key"/></a>
## Retrieving a configuration item by key [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /item/**{item_key}**/|
| Response Content Type | application/json or application/x-yaml |

```bash
# execute the GET operation on the item URI passing its natural key
$ curl \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/item/KEYDEMOM001' 
```

<a name="linking-two-items"/></a>
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
  "meta": { },
  "description": "This is a CMDB item for testing purposes.",
  "role": "connect",
  "start_item_key": "ITEM_ONE_KEY",
  "end_item_key": "ITEM_TWO_KEY",
  "tag": "Test"
}
```

**NOTE**: the *meta* field can contain *any* JSON object.

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
$ curl \
    -X PUT 
    -H 'ContentType: application/json' \
    -H '${AUTH_HEADER}' \
    -d '@link_payload.json' \
    'http://localhost:8080/link/my_link_key/' 
```
<a name="retrieving-a-link"></a>
## Retrieving a link [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /link/**{link_key}**/|
| Response Content Type | application/json or application/x-yaml |

### Example

```bash
# execute the GET operation on the item URI passing the link natural key 
$ curl \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/link/my_link_key/' 
```
