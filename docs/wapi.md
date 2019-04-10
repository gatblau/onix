# Web API 

This section explains how to use Onix Web API. 


<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Web API](#web-api)
  - [Table of Contents (index)](#table-of-contents-index)
  - [Using Swagger (up)](#using-swagger-up)
    - [WAPI User Interface](#wapi-user-interface)
    - [JSON WAPI Docs](#json-wapi-docs)
  - [Identity and Access Management (up)](#identity-and-access-management-up)
  - [Getting Service Information (up)](#getting-service-information-up)
    - [Usage example:](#usage-example)
  - [Creating or updating a configuration item type (up)](#creating-or-updating-a-configuration-item-type-up)
    - [Sample Payload](#sample-payload)
  - [Creating or updating a configuration item (up)](#creating-or-updating-a-configuration-item-up)
    - [Sample Payload](#sample-payload-1)
    - [Payload fields](#payload-fields)
    - [Example](#example)
  - [Retrieving a configuration item by key (up)](#retrieving-a-configuration-item-by-key-up)
  - [Linking two items (up)](#linking-two-items-up)
    - [Sample Payload](#sample-payload-2)
    - [Payload fields](#payload-fields-1)
    - [Example](#example-1)
  - [Retrieving a link (up)](#retrieving-a-link-up)
    - [Example](#example-2)

<a name="using-swagger"></a>
## Using Swagger [(up)](#toc)

Onix uses [Swagger](https://swagger.io/) to document its web API. 

### Swagger UI

When Onix is up and running, the Swagger User Interface can be reached at the following URI:

http://localhost:8080/swagger-ui.html
 
### JSON WAPI Docs

Similarly, a JSON representation of the Web API documentation can be retrieved from the following URI: 
 
http://localhost:8080/v2/api-docs

<a name="identity-and-access-management"></a>
## Access Control [(up)](#toc)

If the Onix Service is configured with '**WAPI_AUTH_MODE=basic**', then a basic authentication token must be passed with every request via an authorisation header.
The following example shows how to obtain a token and pass it to the service:

```bash
# executes the request passing the credentials using the -u otion
curl -u username:password 'http://localhost:8080/itemtype/'
```

To generate a token that can be passed via the http header a generator like [this](https://www.blitter.se/utils/basic-authentication-header-generator/)
can be used. Then, the token can be passed to the API call as follows:

```bash
# executes the request passing the credentials using the -u otion
curl -H 'Authorization: TOKEN_HERE' 'http://localhost:8080/itemtype/'
```

<a name="getting-service-information"></a>
## Getting Service Information [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | GET | 
| Path | /info |
| Response Content Type | application/json |

### Usage example:

```bash
$ curl -u admin:0n1x 'http://localhost:8080/info'
```

<a name="creating-a-configuration-item-type"></a>
## Creating or updating a configuration item type [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | PUT | 
| Path | /itemtype/**{item_type_key}**/|
| Response Content Type | application/json |
 
**NOTE**: this operation is idempotent.

### Sample Payload

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item type is not part of the payload but specified in the URI.

```json
{
  "name": "Item Type 1",
  "description": "This item type is for testing purposes only.",
  "modelKey": "meta_model_1"
}
```

<a name="creating-a-configuration-item"></a>
## Creating or updating a configuration item [(up)](#toc)

| Item  | Value  | 
|---|---|
| Method | PUT | 
| Path | /item/**{item_key}**/|
| Response Content Type | application/json |
 
**NOTE**: this operation is idempotent.

### Sample Payload

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item is not part of the payload but specified in the URI.

```json
{
  "name": "Test Item",
  "description": "This is a CMDB item for testing purposes.",
  "type": "item_type_1",
  "meta": {
    "productId": 1,
    "productName": "A green door",
    "price": 12.50,
    "tags": [ "home", "green" ]
  },
  "tag": ["cmdb", "host", "rhel"],
  "attribute": {
    "WBS" : "012csl",
    "COMPANY" : "ACME"
  },
  "status": 1
}
```
**NOTE**: the *meta* field can contain *any* JSON object.

### Payload fields

The following table describes the fields in the payload and provides some examples of use:

| Field  | Description  | Example |
|---|---|---|
| name | The name associated with the configuration item. | "OCP Demo Master 1" |
| description | A description of the configuration item. | "OpenShift Demo platform first master in zone A."
| type | The natural key of the configuration item type. It must exist as a valid item type. | "HOST" |
| meta | Stores any well-formed json object. This is the primary mechanism to store configuration item information. | { "host":"OCP-DEMO-M-01", "region":"Ireland", "provider":"AWS" } |
| tag | Used for annotating the item for searching. For example, a search can be done by items having the EUROPE tag.| "TEST RELEASE-B EUROPE" |
| status | A number that defines the status of the item, values are arbitrary. | 0 |
| attribute | A set of of key and value pairs. | As per sample payload above. |

### Example

The following example shows how to execute a PUT request to the service using [cURL](https://curl.haxx.se/):

```bash
# execute the PUT operation on the item URI passing a natural key (e.g. KEYDEMOM001) and a payload via json file with contents as per sample above.
$ curl \
    -X PUT \
    -H 'ContentType: application/json' \
    -H 'Authorization: TOKEN_HERE'  
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
    -H 'Authorization: TOKEN_HERE'  
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
  "description": "This is a CMDB item link for testing purposes.",
  "type": "link_type_1",
  "startItemKey": "item_one_key",
  "endItemKey": "item_two_key",
  "tag": ["cmdb", "host", "rhel"],
  "attribute": {
    "WBS" : "012csl",
    "COMPANY" : "ACME"
  },
  "meta": {
    "DeployedBy": "John Wicks",
    "RequestNo": "RN001234"
  }
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
| type | The the natural key of the type of link. | "item_type_x" |
| start_item_key | The key of the item from which the link starts. | A configuration item key. |
| end_item_key | The key of the item to which the link ends. | A configuration item key. |
| attribute | A set of of key and value pairs. | As per sample payload above. |

### Example

The following example shows how to execute a PUT request to the service using [cURL](https://curl.haxx.se/):

```bash
# execute the PUT operation on the item URI passing the link natural key and the payload.json file
$ curl \
    -X PUT 
    -H 'ContentType: application/json' \
    -H 'Authorization: TOKEN_HERE' 
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
    -H 'Authorization: TOKEN_HERE'  
    'http://localhost:8080/link/my_link_key/' 
```
