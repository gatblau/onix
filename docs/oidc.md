# Using OpenId Connect and OAuth 2.0 for authentication and authorisation 
<img src="./pics/ox.png" width="220" height="220" align="right">

- [Using OpenId Connect and OAuth 2.0 for authentication and authorisation](#using-openid-connect-and-oauth-20-for-authentication-and-authorisation)
  - [Introduction](#introduction)
  - [Tested Flows](#tested-flows)
    - [Authorisation Code Flow](#authorisation-code-flow)
    - [Resource Owner Password Flow](#resource-owner-password-flow)
  - [Configuring the Authorisation Server](#configuring-the-authorisation-server)
  - [Configuring the Web API for OIDC](#configuring-the-web-api-for-oidc)
    - [Google OpenId API example](#google-openid-api-example)
    - [OKTA OpenID API example](#okta-openid-api-example)
  - [Authorising users](#authorising-users)
    - [___Onix Scope___](#onix-scope)
    - [___Roles Claim___](#roles-claim)

<a href="introduction"></a>

## Introduction

In the beginning, there were proprietary approaches to working with external identity providers for authentication and authorisation. 

Then SAML (Security Assertion Markup Language) was designed as an open standard using SOAP and XML as its data format. 

Then, there was OAuth and OAuth 2.0 using a RESTful approach to authorisation using JSON as its data format instead.

Today, the preferred way of dealing with “secure delegated access” is OpenID Connect (a.k.a OIDC).

[OpenID Connect](https://openid.net/connect/) is an identity layer running on top of the OAuth 2.0 protocol, that specifies a RESTful HTTP API, using JSON as a data format.

 It allows computing clients to verify the identity of an end-user based on the authentication performed by an authorisation server, as well as to obtain end-user profile information.

<a href="tested-flows"></a>

## Tested Flows

OAuth 2.0 describes a number of grants (“methods”) for a client application to acquire an access token (which represents a user’s permission for the client to access their data) which can be used to authenticate a request to an API endpoint.

The specification describes five grants for acquiring an access token.

Onix has been tested with the two flows described below.

<a href="authorisation-code-flow"></a>

### Authorisation Code Flow

The Authorisation Code Grant Type is probably the most common of the OAuth 2.0 grant types as it is used by both web and native applications to get an access token after a user authorises the application.

The flow has the following steps:

- The application opens a browser to send the user to the OAuth server
- The user sees the authorisation prompt and approves the application's request
- The user is redirected back to the application with an authorisation code in the query string
- The application exchanges the authorisation code for an access token

This flow is seamless to the user when using the Onix Web API for example opening the [Swagger UI](http://localhost:8080/swagger-ui.html).

As the user hits the Swagger User Interface is redirected to the configured authorisation server. Upon entering their credentials and authenticating successfully, the user is redirected to the Swagger UI.

For more information on how thos flow works see [here](https://developer.okta.com/blog/2018/04/10/oauth-authorization-code-grant-type).

<a href="resource-owner-password-flow"></a>

### Resource Owner Password Flow

This flow is typically used by client applications which do not want to use a web browser, for example, when talking to the CMDB using Ansible modules or the Terraform provider.

The flow has the following steps:

- The client application asks the user for their username and password.

- The client then sends a POST request with following body parameters to the authorisation server:

  - __grant_type__ with the value password
  - __client_id__ with the the client’s ID
  - __client_secret__ with the client’s secret
  - __scope__ with a space-delimited list of requested scope permissions.
  - __username__ with the user’s username
  - __password__ with the user’s password

For example, as follows:

```bash
curl  --request POST \
      --url https://dev-447786.okta.com/oauth2/default/v1/token \ 
      --header 'accept: application/json' \
      --header 'Authorization: Basic MG9hbHloMHJ5U1J...' \
      --header 'cache-control: no-cache' \
      --header 'content-type: application/x-www-form-urlencoded' \
      --data 'grant_type=password&username=my_user@email.com&password=dhgcrce..&scope=openid%20onix'
```

__NOTE__: in order to authenticate the request, a [basic access token](https://en.wikipedia.org/wiki/Basic_access_authentication) is passed in the authorization header containing the resource server _client_id_ and _secret_.

- The authorisation server then responds with a JSON object containing the following properties:

  - __token_type__ with the value Bearer
  - __expires_in__ with an integer representing the TTL of the access token
  - __access_token__ the access token itself
  - __refresh_token__ a refresh token that can be used to acquire a new access token when the original expires

```bash
{
  "access_token":"eyJraWQiOiJuWlVs...dfg",
  "token_type":"Bearer",
  "expires_in":3600,
  "scope":"onix openid",
  "id_token":"eyJra...ZOSUZrVGztg"
}
```

The __access_token__ can be decoded using the [jwt.io](https://jwt.io) web site resulting in something like:

```json
{
  "ver": 1,
  "jti": "AT.C6c9uWz4...WfQ...lJ5IAr1LJ...M",
  "iss": "https://dev-447786.okta.com/oauth2/default",
  "aud": "api://default",
  "iat": 15...7858,
  "exp": 15...1458,
  "cid": "0oalyh...gaj356",
  "uid": "00ujnm...zmY356",
  "scp": [
    "onix",
    "openid"
  ],
  "sub": "my_user@email.com",
  "roles": "ADMIN, WRITER"
}
```

__NOTE__: in the example, the token contains a __roles__ claim with two roles ADMIN and WRITER. Onix maps these roles to the privilieges set in its database granting the user access accordingly.

- Finally, the client can then send the bearer token with every request to the Web API, for example to retrieve item types calling the web api as follows:

```bash
curl  --request GET \
      --url http://localhost:8080/itemtype \
      --header 'accept: application/json' \
      --header 'authorization: Bearer eyJraWQiOiJuWlVs...dfg'
```

<a href="configuring-the-authorisation-server"></a>

## Configuring the Authorisation Server

 In order to authenticate using an OpenId Connect provider, a client application and a authorisation server must be created in the Identity provider first.

 As an example, in [Okta](https://developer.okta.com/blog/2017/03/21/spring-boot-oauth#create-an-openid-connect-app-in-okta) do the following:

 To create a new Application:

- Go to Applications
- Add Application
- Select "Native" and click "Next".
- Give it a name
- Tick the "Resource Owner Password" box under **Grant Type Allowed**
- Click done

Now the _client_id_ and _secret_ for the application can be retrieved from the Application General Tab.

The default resource server in Okta can be used.

- Go to API and clieck on Authorisation Servers
- Click on the "default" server link
- Click on "Scopes" and add a new scope called __onix__
- Click on "Claims" and add a new claim called __roles__
- Select "include claim in access token"
- Select "Value Type = Expression"
- Type __String.toUpperCase("admin, writer")__ as value
- Under Include in section select "The following scopes"
- Type __onix__ as the scope

<a href="configuring-the-web-api-for-oidc"></a>

## Configuring the Web API for OIDC

 In addition, the following environment variables should be specified when launching the Onix WAPI application:

 | Variable | Description |
 |---|---|
 | __AUTH_MODE__ | The authentication mode used by the Onix Web API, must be set it to "__oidc__". |
 | __OIDC_ISSUER__ | The name of the token issuer. |
 | __OIDC_JWKURL__ | The URI whwre the issuer has published the JavaScript Web Key certificates used to validate Java Web Tokens. |
 | __OIDC_REDIRECT_URI__ | The resource server endpoint prfroming the login. __Note:__ in the Onix Web API the endpoint is "_/oidc-login_". |
| __OIDC_USER_AUTH_URI__ | The URI of the authorisation server. |
| __OIDC_ACCESS_TOKEN_URI__ | The URI of the token resource. |
| __OIDC_CLIENT_SECRET__ | The client secret used to authenticate with the authorisation server. |
| __OIDC_CLIENT_ID__ | The unique identifier for the client application to be protected. |

__Note:__ client id and secret values are created when setting up a new client application in the OpenId provider.

### Google OpenId API example

The following values apply to a google provider:

| Variable | Value |
|---|---|
| __AUTH_MODE__ | oidc |
| __OIDC_ISSUER__ | accounts.google.com |
| __OIDC_JWKURL__ | https://www.googleapis.com/oauth2/v2/certs |
| __OIDC_REDIRECT_URI__ | http://localhost:8080/oidc-login |
| __OIDC_USER_AUTH_URI__ | https://accounts.google.com/o/oauth2/auth |
| __OIDC_ACCESS_TOKEN_URI__ | https://www.googleapis.com/oauth2/v3/token |
| __OIDC_CLIENT_SECRET__ | 24 character string |
| __OIDC_CLIENT_ID__ | id_string + .apps.googleusercontent.com |

Information on how to create an OpenId Connect application using Google OAuth 2.0 APIs can be found [here](https://developers.google.com/identity/protocols/OpenIDConnect).

### OKTA OpenID API example

The following values apply to an Okta provider:

| Variable | Value |
|---|---|
| __AUTH_MODE__ | oidc |
| __OIDC_ISSUER__ | https://dev-447786.okta.com/oauth2/default |
| __OIDC_JWKURL__ | https://dev-447786.okta.com/oauth2/default/v1/keys |
| __OIDC_REDIRECT_URI__ | http://localhost:8080/oidc-login |
| __OIDC_USER_AUTH_URI__ | https://dev-447786.okta.com/oauth2/default/v1/authorize |
| __OIDC_ACCESS_TOKEN_URI__ | https://dev-447786.okta.com/oauth2/default/v1/token |
| __OIDC_CLIENT_SECRET__ | 40 character string|
| __OIDC_CLIENT_ID__ | 20 character string |

Information on how to create an OpenId Connect application in Okta can be found [here](https://developer.okta.com/blog/2017/03/21/spring-boot-oauth#create-an-openid-connect-app-in-okta).

## Authorising users

OpenID Connect employs OAuth 2.0 access tokens to allow client apps to retrieve user information from the UserInfo endpoint.

Additionally, access tokens can carry claims which are used by the protected resources to determine if the user has access to the resources.

### ___Onix Scope___

An OpenID provider for Onix Web API will have to add a new access token scope called __onix__.

The following picture shows how to add the new scope in Okta:

![Creating an Onix scope](./pics/onix_scope.png "Creating an Onix scope")

### ___Roles Claim___

In order to authorise logged on users, the OAuth 2.0 access token must contain a claim called __roles__ under the __onix__ scope, with a comma separated list of the role names granted to the user.

The following picture shows an example of adding a static roles claim in an Okta OAuth 2.0 access token:

![Adding a role claim](./pics/role_claim.png "Adding role claim to the access token")

Typically, the value of this claim should be dynamically determined based on the user group membership in the backing  directory.

The access token should look like this:

```json
{
  "ver": 1,
  "jti": "AT.0W3...DdQuL-qMGl...wuuK-u6Y...0I",
  "iss": "https://dev-447786.okta.com/oauth2/default",
  "aud": "api://default",
  "iat": 155....369,
  "exp": 155....969,
  "cid": "0oajnv4.....YRtO0356",
  "uid": "00ujnmh.....zmyXY356",
  "scp": [
    "openid",
    "onix"
  ],
  "sub": "user@mail.com",
  "role": "ADMIN, WRITER"
}
```

The roles are then mapped within the Onix WAPI to data model logical partitions via Read, Write and/or Delete privileges according to the model described [here](./rbac.md).
