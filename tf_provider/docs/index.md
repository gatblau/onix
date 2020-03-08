# Onix Terraform Provider  <img src="../docs/pics/ox.png" width="200" height="200" align="right">

This Terraform provider is the default command line interface for managing configuration information in Onix.

For example:

- Creating Models, Item Types, Link Types and Link Rules.
- Creating, updating or destroying items and links. 
- Retrieving configuration information using Data Sources.

## Example Usage

__*Basic Authentication Example*__

In order to authenticate using [Basic Access Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication), there is no need to specify the __auth_mode__ attribute. Only _user_ and _pwd_ are required:

```hcl-terraform
# Basic Authentication
provider "ox" {
  uri       = "http://localhost:8080"
  auth_mode = "basic" # default value if no specified
  user      = "user-name-here"
  pwd       = "user-password-here"
}
```

__*OpenId / OAuth 2.0 Authentication Example*__

If [OpenId Connect / OAuth 2.0](https://openid.net/connect/) is selected as the authentication method, then in addition to the _user_ and _pwd_ attributes, _client_id_ and _secret_ and _token_uri_ are also required:

```hcl-terraform
# OpenId Authentication
provider "ox" {
  uri       = "http://localhost:8080"
  auth_mode = "oidc"
  user      = "user-name-here"
  pwd       = "user-password-here"
  client_id = "application-client-id-here"
  secret    = "application-secret-here"
  token_uri = "uri-of-the-token-endpoint-at-authorisation-server"
}
```

__*No Authentication Example*__


```hcl-terraform
# No Authentication
provider "ox" {
  uri       = "http://localhost:8080"
  auth_mode = "none"
  user      = ""
  pwd       = ""
}
```

__Note__: user & pwd attributes need to be specified but its value is not used if _auth_mode_ is set to _none_.


## Argument Reference

Connection information can be provided by specifying the following attributes of the ox provider as follows:

| Attribute | Description | Example |
|---|---|---|
| __uri__| The URI of the Onix Web API | http://localhost:8080 |
| __auth_mode__ | Defines the method used by the provider to authenticate with the Onix Web API. If not specified, it defaults to __basic__ (basic authentication). Other possible value as __none__ or __oidc__ (OpenId Connect). | basic |
| __client_id__ | The public identifier for the Onix Web API defined by the OAUth 2.0 server. It is only required if _auth_mode_ is set to _oidc_. | 2Idlxf0ryAGOd3gaj938 |
| __secret__ | A secret known only to the application and the authorisation server. It is only required if _auth_mode_ is set to _oidc_. | Hl5_V_lbhLQHol47f5is6YErs7pHKP3OP3oEf7H3 |
| __token_uri__ | The OAuth 2.0 server endpoint where the ox provider exchanges the user credentials, client ID and client secret, for an access token. It is only required if _auth_mode_ is set to _oidc_. | https://dev-1234.okta.com/oauth2/default/v1/token |
 | __user__ | A unique sequence of characters used to identify a user of the Onix Web API. A typical value could be the user email address defined in the OAuth Server. | user@email.com |
 | __pwd__ | A secret word supplied by the user in order to gain access to the Onix Web API. | 0n1x_440d4f6f |
  
__NOTE__: The __auth_mode__ attribute must match the value used by the Onix Web API. For example, if the Onix Web API is set to use __auth_mode=oidc__ then the terraform provider must be set to use the same __auth_mode__, otherwise the authentication will fail.

