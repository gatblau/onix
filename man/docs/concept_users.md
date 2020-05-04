---
id: concept_users
title: Users
---
Anyone who wants to use Onix requires an account, so that they can be identified and allocated privileges to create,
update, list and delete configuration data.

Accounts in Onix are represented by *Users*. 

Depending on the authentication method selected, **users** can be stored in:

1. The Onix database, when using [Basic Access Authentication](https://en.wikipedia.org/wiki/Basic_access_authentication).
2. An external system, when using [OpenId Connect](https://en.wikipedia.org/wiki/OpenID_Connect).

## Local Users

Onix allows creating, updating, listing and deleting local users via its Web API. When doing so, user information 
is stored in its database. 

Local users have the following data attributes:

| attribute | description | 
|---|---|
| key | the user unique natural key used to retrieve and update the user data |
| name | the unique username |
| email | the unique user email address |
| pwd | the user password |

Storing users locally is useful in cases where integration with an organisation's identity and access management system 
is not desired. It is provided as a way to get up and running easily whilst giving complete access control to configuration
data.

### User registration

Typically, onboarding users securely can be a complex task. It involves various steps to activate an account and ensure the 
user is legitimate. In order to address this challenge, Onix provides an email notification system to facilitate the
registration process as follows:

| actor | action |
|---|---|
| **System Admin** | creates a User in the system using an email address |
| **System** | emails the user notifying them they need to reset their password |
| **User** | accesses the system and request a password reset token |
| **System** | emails the password reset token to the user |
| **User** | enters their password using the reset token |
| **System** | validates the token and changes the password |
| **System** | emails the user to notify them their password has been changed |

:::note

*If the above process does not fit the organisation requirements, then it is recommended the use of external users.*

External user integration allows Onix to reuse any users and roles that are set up as part of the organisation's 
procedures, and is explained in more detail in the *external users* section below.

:::

### Resetting Passwords

Resetting passwords for local users is done in the same way as in the user's first registration process.

Any user who has forgotten their password, can request a password reset token and then use it to change their password 
as follows:

| actor | action |
|---|---|
| **User** | accesses the system and request a password reset token |
| **System** | emails the password reset token to the user |
| **User** | enters their password using the reset token |
| **System** | validates the token and changes the password |
| **System** | emails the user to notify them their password has been changed |

### Memberships

In order to get access to data resources, users require the right privileges. 

Privileges are grouped into Roles and Memberships link a user with one or more roles.

Therefore, in order to access resources, a user must have one or more memberships to one or more roles.

The key attributes held by a membership are:

| attribute | description | 
|---|---|
| key | the membership unique natural key used to retrieve and update the membership data |
| user | the user the membership is for |
| role | the role granted to the user via this membership |

### Default users

Onix ships out of the box with three default local users as follows:

| user | privileges |
|---|---|
| admin | a super administrator who can read and write data and models |
| reader | a read only user who can only read data |
| writer | a read and write user who can read and write data |

:::note

The above users are members of the ADMIN, READER and WRITER roles respectively. The extent of their privileges is 
discussed in the roles section.

:::

## External Users

External users are not stored in Onix. The idea is that Onix only needs to know the username and the roles the user is in.

When the authentication mode is set to OIDC (OpenId Connect), Onix authenticates users based on an OpenId token which 
carries a list of roles to which the user is entitled.

In this way, and in contrast to local users, users and memberships are not stored in Onix but are part of an external Identity & Access Management System.

The only prerequisite is that the organisation has an OpenID server that can provide tokens with membership claims (i.e.
 the list of roles for the user).
 
:::important

The role names in the OpenId token must match roles in the Onix database.

:::

:::tip

Not storing users and memberships in Onix allows an organisation to leverage existing onboarding processes and 
security procedures that comply with their specific requirements.

:::
