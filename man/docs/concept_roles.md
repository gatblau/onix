---
id: concept_roles
title: Roles
---
import useBaseUrl from '@docusaurus/useBaseUrl';

A role is a group of [privileges](/onix/docs/concept_roles#privileges) that can be assigned to [users](/onix/docs/concept_users) to *create*, *list* and/or *delete* configuration data within 
a [partition](/onix/docs/concept_roles#partitions).
Partitions and privileges are explained below.

When a user logs in, Onix obtains a list of roles that have been assigned to the user and pass them to its data service.

The data service *only* creates, lists or deletes data to which the user has been granted privileges - through their role(s).

In order to apply privileges to *only certain portions of the data stored in Onix*, it is necessary to partition the information    as explained below.

## Partitions

Partitions are a logical way to divide the Onix data into areas so that privileges can be given to access such areas 
independently. 

For example, suppose that a specific group within an organisation wants to be able to create and manage configuration 
data that is not supposed to be seen or managed by another group within the organisation. In this case, a partition is 
created for such group, and privileges are given to a role to *Create, Read and / or Delete* data in the partition.

Once users are [members](/onix/docs/concept_users#memberships) of the role, they inherit those privileges.

### Default partitions

In order to facilitate getting started and to provide the simplest access controls possible, Onix ships out of the box with
 two default partitions:

| key | partition | description |
|---|---|
| REF | The default *reference data* partition. By default, any reference data in Onix is put into the REF partition. Reference data is any data that is part of configuration models (i.e. models, item types, link types and link rules). |
| INS | The default *instance data* partition. By default, any instance data in Onix is put into the INS partition. Instance data is comprised by items and links. Items and links are explained in later sections. |

The following figure shows the two types of partitions (i.e. reference and instance) as well as the data entities to which
they apply:

<img alt="Partitions" src={useBaseUrl('img/concept_partitions.png')} />

## Administrator roles

In addition to the default partitions, administrators can create and update partitions.

Administrator roles can create or update partitions. 
They are differentiated from ordinary roles by their *level* which is either 1 or 2.

The following role levels are possible:

| role level | type | description |
|:-:|---|---|
| 0 | ordinary user | cannot manage partitions |
| 1 | local administrator | can only manage the partitions they create |
| 2 | super administrator | can manage all partitions irrespective of who have created them |

## Privileges

A privilege is a permission for a user to perform an operation on a piece of configuration data.

There are three types of operations for which a User can be granted privileged:

| privilege operation  | description |
|:-:|---|
| CREATE | allows users to create or update data. Onix operations are [idempotent](https://en.wikipedia.org/wiki/Idempotence). No distinction is made between create and update operations. |
| READ | allows users to list data. |
| DELETE | allows users to delete data. |

<img alt="From Users to Partitions" src={useBaseUrl('img/concept_users_to_parts.png')} />

