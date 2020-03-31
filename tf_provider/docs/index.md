# Onix Terraform Provider  <img src="../../docs/pics/ox.png" width="200" height="200" align="right">

This Terraform provider is the default command line interface for managing configuration information in [Onix](https://onix.gatblau.org).

For example:

- Creating Models, Item Types, Link Types and Link Rules.
- Creating, updating or destroying items and links.
- Retrieving configuration information using Data Sources.

## Provider Configuration Argument Reference

The provider can be configured via attributes in the provider section of the Terraform file or via environment variables.

It is recommended that for production use, environment variables are used to avoid placing credentials in the Terraform file and inadvertedly commiting them into a source control system.

The following table shows the provider configuration attributes:

| provider<br>attribute | environment<br>variable | use | description | example |
|---|---|---|---|---|
| `uri` | `TF_PROVIDER_OX_URI` | *required* | *The URI of Onix Web API where the provider will connect.* | `http://localhost:8080` |
| `auth_mode` | `TF_PROVIDER_OX_AUTH_MODE` | optional | *Defines the method used by the provider to authenticate with the Onix Web API. If not specified, it defaults to __basic__ (basic authentication). Other possible value as __none__ or __oidc__ (OpenId Connect).* | `basic` |
| `user` | `TF_PROVIDER_OX_USER` | *required*<br>(if `basic`) | *A unique sequence of characters used to identify a user of the Onix Web API.* | `admin`, `reader` or `writer` |
| `pwd` | `TF_PROVIDER_OX_PWD` | *required*<br>(if `basic`) | *A secret word supplied by the user in order to gain access to the Onix Web API.* | `pwd0012asx!` |
| `token_uri` | `TF_PROVIDER_OX_TOKEN_URI` | *required*<br>(if `oidc`) | *The OAuth 2.0 server endpoint where the ox provider exchanges the user credentials, client ID and client secret, for an access token. It is only required if _auth_mode_ is set to _oidc_.* | `https://token-server.com/oauth2/default/v1/token` |
| `app_client_id` | `TF_PROVIDER_OX_APP_CLIENT_ID` | *required*<br>(if `oidc`) | *The public identifier for the Onix Web API defined by the OAUth 2.0 server. It is only required if _auth_mode_ is set to _oidc_.* | `character string of lenght determined by implementation` |
| `app_secret` | `TF_PROVIDER_OX_APP_SECRET` | *required*<br>(if `oidc`) | *A secret known only to the application and the authorisation server. It is only required if _auth_mode_ is set to _oidc_.* | `character string of lenght determined by implementation` |

__NOTE__: The __auth_mode__ attribute must match the value used by the Onix Web API. For example, if the Onix Web API is set to use __auth_mode=oidc__ then the terraform provider must be set to use the same __auth_mode__, otherwise the authentication will fail.

## Example Usage

### Basic Authentication Example (configuration in terraform file)

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

### Basic Authentication Example (configuration in environment variables)

As a good practice, it is recommended that any credentials are not specified in the terraform file, but instead, they are provided via environment variables.

For example:

On the Terraform file add the section to select the provider:

```hcl-terraform
provider "ox" {
}
```

On the command line, specify the provider configuration through environment variables:

```bash
$ TF_PROVIDER_OX_URI=http://localhost:8080 \
  TF_PROVIDER_OX_AUTH_MODE=basic \
  TF_PROVIDER_OX_USER=admin \
  TF_PROVIDER_OX_PWD=0n1x \
  terraform apply
```

### OpenId / OAuth 2.0 Authentication Example

If [OpenId Connect / OAuth 2.0](https://openid.net/connect/) is selected as the authentication method, then in addition to the _user_ and _pwd_ attributes, _client_id_ and _secret_ and _token_uri_ are also required:

```hcl-terraform
# OpenId Authentication
provider "ox" {
  uri           = "http://localhost:8080"
  auth_mode     = "oidc"
  user          = "user-name-here"
  pwd           = "user-password-here"
  app_client_id = "application-client-id-here"
  app_secret    = "application-secret-here"
  token_uri     = "uri-of-the-token-endpoint-at-authorisation-server"
}
```

### No Authentication Example

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

## Terraform Resources

A list of available resources can be found [here](resources/index.md).

## Terraform Data Sources

A list of available data sources can be found [here](datasources/index.md).
