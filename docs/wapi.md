# Web API <img src="./pics/ox.png" width="160" height="160" align="right">

This section explains how to use Onix Web API. 

<a name="toc"></a>
## Quick Index [(back)](./../readme.md)

### _Good to know_

| _Area_ | _Description_ |
|---|---|
| [Using Swagger](#using-swagger) | How to access online Web API documents? |
| [Access Control](#access-control)| How to authenticate and authorise users? |
| [Getting WAPI information](#getting-wapi-info)| How to get the Web API version online? |
| [HTTP return codes](#http-return-codes)| Which codes Web API return? |
| [HTTP Result](#http-result)| What data is returned by the service when resources are created, updated or deleted? |
| [Concurrency Management](#concurrency-management)| How to use the Web API in concurrent user scenarios? |
| [Automation Clients](#automation-clients)| How do I apply this documentation to Ansible and Terraform clients? |
| [Idempotence](#idempotence)| What happens if I make repeated requests to the same endpoint with the same payload? |
| [Change History](#change-history)| Is every change made recorded? |
| [Versioning](#versioning)| Can I version CMDB data? |

### _Reference data: working with models_

| _Area_ | _Description_ |
|---|---|
| [Models](#models)| How to create a new or update or delete an existing model definition? |
| [Item types](#item-types)| How to create a new or update or delete an existing item definition? |
| [Link types](#link-types)| How to create a new or update or delete an existing link definition? |
| [Link rules](#link-rules)| How to create a new or update or delete an existing link rule? |

### _Instance data: feeding and querying the database_

| _Area_ | _Description_ |
|---|---|
| [Items](#items)| How to create a new or update or delete an existing item type? |
| [Links](#links)| How to create a new or update or delete an existing link type? |


<a name="using-swagger"></a>
## Using Swagger [(up)](#toc)

Onix uses [Swagger](https://swagger.io/) to document its web API. 

### Swagger UI

When Onix is up and running, the Swagger User Interface can be reached at the following URI:

http://localhost:8080/swagger-ui.html
 
### JSON WAPI Docs

Similarly, a JSON representation of the Web API documentation can be retrieved from the following URI: 
 
http://localhost:8080/v2/api-docs

<a name="access-control"></a>
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

### Authentication Modes

| _Mode_ | _Description_ | _Setting_ |
|---|---|---|
| None | No authentication is required. This is mainly used for development activities where authentication is not required. | AUTH_MODE=none|
| Basic Auth| Basic authentication distinct reader, writer and admin credentials. | AUTH_MODE=basic |
|OIDC[*] | OpenID Connect supporting identity and access tokens. | AUTH_MODE=oidc |

[*] This is currently work in progress and subject of the next release.

<a name="getting-wapi-info"></a>

## Getting WAPI Information [(up)](#toc)

The root of the service displays service information.

| Item  | Value  |
|---|---|
| Method | GET |
| Path | / |
| Response Content Type | application/json |

### Usage example

```bash
$ curl -u admin:0n1x 'http://localhost:8080'
```

### Return values

| _Attribute_ | _Description_ |
|---|---|
| description | name of the service |
| version | version of the release |

<a name="http-return-codes"></a>

## HTTP Return Codes [(up)](#toc)

For any resource acting on CMDB entities, the following codes are returned by the HTTP requests:

| _Code_ | _Method_ | _Description_ |
|---|---|---|
| 200 | GET, PUT (U, L, N), DELETE | Successful HTTP request |
| 201 | PUT (I) | Successful HTTP request |
| 401 | PUT, GET, DELETE | Unauthorised HTTP request |
| 404 | PUT, GET, DELETE | Resource not found |
| 500 | PUT, GET, DELETE | Internal server error |

For PUT, DELETE and GET HTTP methods, the following operations are available:

| _Code_ | _Operation_ | _Description_ |
|---|---|---|
| __I__ | Insert | entity inserted |
| __U__ | Update | entity updated |
| __D__ | Delete | entity deleted |
| __L__ | Lock | no action taken, version on the client does not match the version on the server
| __N__ | None | no action taken, client and server data match |

<a name="http-result"></a>

## HTTP Result [(up)](#toc)

In the case of PUT and DELETE HTTP methods, a Result object is returned containing information about the request process as follows:

| _Attribute_ | _Description_ |
|---|---|
| __ref__ | a reference containing the entity type and natural key |
| __operation__ | a character indicating the type of operation executed; namely I, U, D, L or N as described in the table above.
| __changed__ | true if the entity was created, updated or deleted |
| __error__ | true if there was a problem processing the request |
| __message__ | an error message in case _error_ is true |

<a name="concurrency-management"></a>

## Concurrency Management [(up)](#toc)

Every time a resource is updated, an incremental version number is automatically assigned to the resource.

All resource requests can accept a version number when an HTTP PUT method is requested. The version number enables [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control).

If no version number is specified in the PUT request, concurrency is disabled. However, if a value is specified, the following outcomes are possible:

1. __Insert or Update__: an insert/update occurs if the client version number matches the server version number.
2. __Optimistic Lock__: no action is taken if the client version number is behind the server version number. The response contains an __L__ operation which means that some other client has updated the state of the resource since the last time a copy was retrieved. This feature is helpful for user interfaces updating resources where more than one client could be acting on the same resource at the same time. The client has the option to: a) override the server by sending the request again without a version number; or b) refreshing the client with the new data from the server.

<a name="automation-clients"></a>

## Automation Clients [(up)](#toc)

Onix integrates with [Ansible](https://www.ansible.com/) (via Ansible Modules) and [Terraform](https://www.terraform.io/) (via a Terraform Provider).

Each of the Web API resources maps directly to an Ansible Module or a Terraform Resource/Data Source. Therefore, in order to avoid maintaining multiple document sets, this document should be used to determine which attributes are available for each resource within both, Ansible and Terraform (TF).

The following table provides a convention for translating Web API resources into Ansible module names, Terraform resource names and Terraform data source names:

| Web API Resource | Ansible Module | TF Resource | TF Data Source |
|---|---|---|---|
| model |ox_model| ox_model | ox_model_data |
| itemtype |ox_item_type |ox_item_type | ox_item_type_data |
| linktype | ox_link_type | ox_link_type | ox_link_type_data |
| linkrule | ox_link_rule | ox_link_rule | ox_link_rule_data |
| item | ox_item | ox_item | ox_item_data |
| link | ox_link |ox_link |ox_link_data |

<a name="idempotence"></a>

## Idempotence [(up)](#toc)

All HTTP resources in the Web API are [idempotent](https://en.wikipedia.org/wiki/Idempotence). 

Resources are uniquely identified by a [natural key](https://en.wikipedia.org/wiki/Natural_key). The natural key is a string with a value that should be decided using a naming convention that is chosen by the implementer.

Using a natural key, for example, in HTTP PUT methods, allow the client to avoid having to think if the HTTP resource has to be created or updated. The request ensures that the resource information is in the CMDB regardless of how many times the method is invoqued.

When a PUT method is executed for the first time, the resource is created. Any subsequent calls to the same method will result in updates if the payload has changed, or no action will be taken if the payload is the same.

<a name="change-history"></a>

## Change History [(up)](#toc)

Every time a CMDB entity is created, updated or deleted, their data get written into change tables recording the type of operation (i.e. **I**nsert, **U**pdate or **D**elete), the date of the change and the user that made the change.

Preserving the whole history of changes allows to retrieve any configuration at any point in time.

<a name="versioning"></a>

## Versioning [(up)](#toc)

Every entity in the CMDB has a version number starting with 1, which increments automatically when changes are made. 

In order to retrieve a configuration, that is a set of items and links at a specific point in time,  a **tag** needs to be created. The tag carries the information of which items and links where related to a specified item at a specified point in time.

Creating a tag in the CMDB is equivalent to creating a tag in Git.

<a name="models"></a>

## Models [(up)](#toc)

In order to create items, it is first necessary to create a model, that is a set of item and link definitions (i.e. item types and link types).

Item and Link Types have to belong to one and only one model.
A model can be created as described below.

### Create or Update

To create or update a model use the PUT HTTP method.

#### Request attributes

| _Item_  | _Value_ |
|---|---|
| Method | PUT |
| Path | /model/**{model_key}**|
| Response Content Type | application/json |
| Authentication Header | basic authentication token |

#### Request payload

| _Attribute_ | _Description_ | _Example_ | _Mandatory_|
|---|---|---|---|
| __name__ | the human readable name of the model | "AWS EC2  Model"| yes (unique) |
| __description__ | the model description | "Definitions for AWS Elastic Compute Cloud items and their relationships" | no |
| __version__ | the version of the model for concurrency management purposes. | 34 | no |

#### Usage

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item type is not part of the payload but specified in the URI.
```bash
$ curl \
    -X PUT \
    -H 'ContentType: application/json' \
    -H '${AUTH_HEADER}' \
    -d '@model_payload.json' \
    'http://localhost:8080/model/awsec2' 
```

__model_payload.json__:

```json
{
  "name": "AWS EC2 Model",
  "description": "Definitions for AWS Elastic Compute Cloud items and their relationships."
}
```

__result__:

```json
{
  "ref": "model:awsec2",
  "changed": "true",
  "error": "false",
  "message": "", 
  "operation": "I", 
}
```

### Delete

In order to delete a model the following command can be executed:

```bash
$ curl \
    -X DELETE \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/model/awsec2' 
```

__NOTE__: deleting a Model, forces the deletion of the Link Types and Item Types related to that Model. Item and Links have to be deleted before a model can be deleted.

### Query

The following model queries are available:

| _Query_ | _Description_ |
|---|---|
| __GET /model/{model_key}__ | Retrieve the model for the specified natural key. |
| __GET /model/{model_key}/data__| Retrieve a list of Item Types and Link Types that comprise the specified model. |
| __GET /models__ | Retrieve a list of all models in the system. |

--------------------------------

<a name="item-types"></a>

## Item Types [(up)](#toc)

Item Types provide the definition of items. Items can record data in two main ways:

- Using a JSON object (meta field)
- Using Key/Value pairs (attribute field)

Item types can be used to apply specific validation to the data recorded by a configuration item by defining a JSON schema or key/value required/mandatory constraints.

### Create or Update

To create or update an Item Type use the PUT HTTP method.

#### Request attributes

| _Item_  | _Value_ |
|---|---|
| Method | PUT |
| Path | /itemtype/**{item_type_key}**|
| Response Content Type | application/json |
| Authentication Header | basic authentication token |

#### Request payload

| _Attribute_ | _Description_ | _Example_ | _Mandatory_|
|---|---|---|---|
| __name__ | the human readable name of the item type | "AWS  Instance"| yes (unique) |
| __description__ | the item type description | "A Virtual Machine running on Aws EC2" | no |
| __version__ | the version of the model for concurrency management purposes. | 21 | no |

#### Usage

The PUT request requires a payload in JSON format as the one shown below.
Note that the natural key for the configuration item type is not part of the payload but specified in the URI.
```bash
$ curl \
    -X PUT \
    -H 'ContentType: application/json' \
    -H '${AUTH_HEADER}' \
    -d '@model_payload.json' \
    'http://localhost:8080/model/awsec2' 
```

__model_payload.json__:

```json
{
  "name": "AWS EC2 Model",
  "description": "Definitions for AWS Elastic Compute Cloud items and their relationships."
}
```

__result__:

```json
{
  "ref": "model:awsec2",
  "changed": "true",
  "error": "false",
  "message": "", 
  "operation": "I", 
}
```

### Delete

In order to delete a model the following command can be executed:

```bash
$ curl \
    -X DELETE \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/model/awsec2' 
```

__NOTE__: deleting a Model, forces the deletion of the Link Types and Item Types related to that Model. Item and Links have to be deleted before a model can be deleted.

### Query

The following model queries are available:

| _Query_ | _Description_ |
|---|---|
| __GET /model/{model_key}__ | Retrieve the model for the specified natural key. |
| __GET /model/{model_key}/data__| Retrieve a list of Item Types and Link Types that comprise the specified model. |
| __GET /models__ | Retrieve a list of all models in the system. |

--------------------------------