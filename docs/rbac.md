# Role Based Access Controls <img src="./pics/ox.png" width="160" height="160" align="right">

  - [Semantic Model](#semantic-model)
  - [Default partitions](#default-partitions)
  - [Default Roles](#default-roles)
  - [Authentication Modes](#authentication-modes)
  - [Controlling access witn OpenId Connect](#controlling-access)

<a href="semantic-model"></a>
## Semantic Model

In order to restrict access to certain portions of the CMDB, Onix divides the data into logical partitions.

Users, logged under a specific role, can then be granted privilege to create/update, read and or delete data in a logical partition.

The following diagram shows how access control works:

![Onix Data Model](./pics/semantic_rbac_model.png "Role Based Access Control").

Partitions can be assigned to either:

- A __model__ and all its related item types, link types and link rules.
- An __item__ and any links associated with it.

<a href="default-partitions"></a>
## Default partitions

By default, Onix provides two partitions as follows:

| Partition | Description |
|---|---|
| __REF__ | The reference data default partition. Any model information is added to this partition by default, unless the model data is explicitly assigned to a different partition. |
| __INS__ | The instance data default partition. Any item and related links are added to this partition by default, unless configuration data explicitly assigned to a different partition. |

<a href="default-roles"></a>
## Default Roles

By default, Onix provides three roles as follows:

| Role | Description |
|---|---|
| __ADMIN__ | Has default privilege to create/update, read and delete on both the REF and INS partitions. |
| __READER__ | Has default privilege to read on both the REF and the INS partitions. |
| __WRITER__ | Has default privilege to read on the REF partition and to create/update, read and delete on the INS partition. |

In other words:

| Role | Partition | Create? | Read? | Delete? |
|---|:--:|:--:|:--:|:--:|
| __ADMIN__  | REF | Y | Y | Y |
| __ADMIN__  | INS | Y | Y | Y |
| __WRITER__ | REF | N | Y | N |
| __WRITER__ | INS | Y | Y | Y |
| __READER__ | REF | N | Y | N |
| __READER__ | INS | N | Y | N |

<a href="authentication-modes"></a>
## Authentication Modes

Onix supports three authentication modes, based on the value of the __AUTH_MODE__ environment variable as follows:

| AUTH_MODE | Description | RBAC Behaviour |
|---|---|---|
| __none__ | No authentication is required. Intended only for development. | Role ADMIN is applied to the user by default. |
| __basic__ | Basic authentication is required. | Three fix set of credentials can be configured at the service level mapping to ADMIN, WRITER and READER roles above. |
| [__oidc__](./oidc.md) | [OpenId Connect authentication](./oidc.md) is required. | The token presented to the service must have a claim called __role__ containing the name to be mapped to the role defined in the Onix data model. This role in turn defines the privileges the user has on one or more logical partitions. |

<a href="controlling-access"></a>
## Controlling access to data with [OpenId Connect](./oidc.md)

Suppose you want to restrict access to read part of the data model to a group of users. For instance, a particular project team needs to read and write data for the project but users of other projects should not have access to this data.

A logical partition is created for example with the name of the project. Any data the project creates in the CMDB is set to use such partition.

In parallel, two roles are created, namely *Project_Reader* and *Project_Writer*.

Privileges are also created for each role. The *Project_Reader* is granted read access to the *Project* partition whereas the *Project_Writer* is granted create, read and delete rights to the *Project* partition.

The OpenId provider fronts a user directory where two groups exist for the project, namely *Project_Reader* and *Project_Writer*.

When the user access the Onix service, the service request the provider to authenticate the user and issue a token.

The provider verifies the user credentials and if successful, adds a claim named __role__ to the token. The value of the claim is based on the group membership of the user in the directory. For example, if the user is in the *Project_Writer* role, then the claim role contains the *Project_Writer* value.

The token is passed back to Onix, which extracts the claim role and match it against the role defined in the database. Any data operation requested by the user is subjected to verification of its role privileges.
If the role does not have a privilege to execute the requested data operation, an exception is raised.






