# Identity and Access Management 

Onix CMDB by default uses [Keycloak](https://www.keycloak.org/) as its Identity and Access Management (IDAM) solution.

<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Setting up Keycloak](#setting-up-keycloak)
- [Configuring Keycloak](#configuring-keycloak)
- [Onix Roles](#onix-roles)
- [Requesting an Access Token](#requesting-an-access-token)
- [WAPI configuration variables for IDAM](#wapi-configuration-variables-for-idam)
- [Issuing calls to the WAPI](#issuing-calls-to-the-wapi)

<a name="setting-up-keycloak"/></a>
## Setting up Keycloak [(up)](#toc)

To understand how to set up Keycloak see [here](https://www.keycloak.org/docs/latest/server_installation/index.html). 

If using the [docker-compose.yml](../install/container/docker-compose.yml) file to start all relevant containers, then Keycloak will be launched as a container with an embedded H2 in memory database. 

If not, and for testing purposes, a simple docker container running Keycloak can be created as follows:

```bash
$ docker run --name idam -p 8081:8080 -d -e KEYCLOAK_USER=admin -e KEYCLOAK_PASSWORD=admin jboss/keycloak
```

<a name="configuring-keycloak"></a>
## Configuring Keycloak [(up)](#toc)

To configure Keycloak for Onix CMDB do the following:
 
- Create a docker container running Keycloak by running the [run.sh](../install/container/idam/run.sh) script. This step is not needed if Keycloak is started as part of 
- Launch Keycloak on http://localhost:8081/ 
- Log using admin/admin
- Import the following [realm file](../install/container/idam/onix-realm.json).

Alternatively, follow the manual steps below:
- Create a realm named "onix"
- Create a new client within the realm called "onix-cmdb"
- Set the client's "Valid Redirect URIs" to "http://localhost:8080/*" (assuming the web API is listening to the localhost on port 8080, otherwise change it to the correct regex.)
- Create a new role called "admin"
- Create a new role called "user"
- Create a new user called "onix"
- Set the user password
- Within users, role mappings, add the role "user" or "admin" to the new user.

## Onix Roles [(up)](#toc)

|Role | Description |
|---|---|
| **user** | Has read only access to the resources under the /info path. This role is intended to be used to query CMDB information. |
| **admin** | Has read and write access to resources under the root / path. This role is required to login using the Ansible module [onix_login](../ansible/library/onix_login.py). It should be used to create and update items and links. |


<a name="requesting-an-access-token"></a>
## Requesting an Access Token [(up)](#toc)

Now that keycloak is running and configured, try to get a bearer token running the command below:

```bash
# gets a token from the Auth server
# NOTE: replace the payload attributes depending on the configuration of the Auth server
$ TOKEN=`curl -d "client_id=onix-cmdb" -d "username=onix" -d "password=onix" -d "grant_type=password" "http://localhost:8081/auth/realms/onix/protocol/openid-connect/token"`

$ echo $TOKEN

# constructs an authorization header
# NOTE: check the access_token substrings below is OK for your case
$ AUTH_HEADER='Authorization: bearer '${TOKEN:17:1135} 

$ echo $AUTH_HEADER
```
Where the payload contains the following parameters:

| Parameter | Value  | Description |
|---|---|---|
| **client_id** | onix-cmdb | The identifier for the Onix CMDB application within Keycloak. |
| **username** | onix | The name of the user requesting the token. |
| **password** | onix | The user password. |
| **grant_type** | password | The type of credentials supplied. |

The response should be similar to the json file below:

```json
{ 
  "access_token":"eyJhbGciOi...",
  "expires_in":300,
  "refresh_expires_in":1800,
  "refresh_token":"eyJhbGciO...",
  "token_type":"bearer",
  "not-before-policy":0,
  "session_state":"850e8b86-e759-4a26-a082-ba2a03472cb0"
}
```
<a name="wapi-configuration-variables-for-idam"></a>
## WAPI configuration variables for IDAM [(up)](#toc)

Onix Service jar file can be further configured with the following environment variables:

| Variable | Value  | Default |
|---|---|---|
|**AUTH_SERVER_URL**| The location of Keycloak auth endpoint. | http://localhost:8081/auth |
|**AUTH_REALM**| The application security realm. | onix |
|**AUTH_RESX**| The application to be authenticated (client_id). | onix-cmdb |
|**AUTH_ENABLED**| Whether IDAM is enabled. IDAM is disabled by default, to facilitate testing and getting started with the service. To enable it, set the variable to true. | false | 

<a name="issuing-calls-to-the-wapi"></a>
## Issuing calls to the WAPI 

Once in possession of the bearer token (i.e. access_token) a call to the WAPI endpoint can be issued passing the token via an authorisation header as follows:

```bash
$ curl \
    -H '${AUTH_HEADER}' \
    'http://localhost:8080/itemType/'
```