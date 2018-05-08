# CMDB Model

<a name="toc"></a>
## Table of Contents [(index)](./../readme.md)

- [Overview](#overview)
- [Item Entity](#item)
- [Item Type Entity](#item-type)
- [Link Entity](#link)
- [Dimension Entity](#dimension)
- [Auditing changes](auditing-changes)

<a name="overview"></a>
## Overview [(up)](#toc)

In order to store configuration data Onix uses a simple relational model as shown below:

![Onix Data Model](./pics/onix_model.png "Onix Data Model")

At the moment, the database schema is implemented in PostgresSQL, although in theory, other databases could be used as the Onix Web API uses [Hibernate](https://github.com/hibernate/hibernate-orm) to connect to the database.

<a name="item"></a>
## Item Entity [(up)](#toc)

The Item entity stores information for configuration items. 

Items can be anything that needs to be recorded for example, Virtual and Physical Servers, Middleware, People, etc.

For user purposes, an Item is uniquely identified using a natural key, that is, a unique key formed of attributes which already exist in the configuration domain.
For instance, a Virtual Machine could be identified with a code like "VM-EUR-DC01-ABC". The natural key should never change for that item as it is used to identify the item and execute subsequent updates once the item has been created.

For database relational purposes, an item has a surrogate key used by [Hibernate](https://github.com/hibernate/hibernate-orm) to manage the entities.

### Attributes

|Name | Description | Data Type|
|---|---|---|
|id | Surrogate key. | bigint |
|key| Natural key. | ﻿character varying(100)|
|name| A user friendly name assigned to the item for displaying purposes. | character varying(200) |
|description| The description of the item | text |
|meta| A JSON object of any structure, containing all the specific information for a given configuration item. | json |
|tag| A text field to support semantic tagging. |﻿character varying(300) |
|created| The date and time the item was created. | ﻿timestamp(6) with time zone|
|updated| The date and time the item was last updated. | ﻿timestamp(6) with time zone |
|item_type_id| The type of configuration item. | integer |

<a name="item-type"></a>
## Item Type Entity [(up)](#toc)

The Item Type entity contains the definition of the various types of configuration items.

In order to use the CMDB, item types should be created according to the needs of the particular configuration domain.
Some types are reserved for application use and therefore they should not be deleted (they are marked is a custom=false attribute).

### Attributes

|Name | Description | Data Type|
|---|---|---|
|id | Surrogate key. | bigint |
|name| A user friendly name for the type. | ﻿character varying(200)|
|description| The description of the item type. | ﻿character varying(500) |
|custom| A flag used by the WAPI to discern between system and custom types. | boolean |

<a name="link"></a>
## Link Entity [(up)](#toc)

The Link entity is an association between two items. 

Links allow to create relationships between items for reporting purposes and can optionally store link specific metadata in the association.

Like Items, Links have to be uniquely identified using natural keys. 
 
### Attributes

|Name | Description | Data Type|
|---|---|---|
|id | Surrogate key. | bigint |
|key| Natural key. | ﻿character varying(100)|
|role| A friendly description of the action performed by the link, for example "connects", "hosts", etc. | ﻿character varying(200) |
|description| The description of the link. | text |
|meta| A JSON object of any structure, containing all the specific information for a given link. | json |
|created| The date and time the link was created. | timestamp(6) with time zone |
|updated| The date and time the link was last updated. | ﻿timestamp(6) with time zone |
|start_item_id| The surrogate key of one of the items associated by the link. | bigint |
|end_item_id| The surrogate key of the other item associated by the link. | bigint |

<a name="dimension"></a>
## Dimension Entity [(up)](#toc)

The dimension is a mechanism to attach an indefinite number of key/value pairs for reporting purposes.

For example, Items might be associated with projects so it would be beneficial to know the project work break down structure (WBS) code for the item.

Then, a reporting dimension could be added to the item as follows: key="WBS" and value="OP011.896WIE".

### Attributes

|Name | Description | Data Type|
|---|---|---|
|id | Surrogate key. | bigint |
|item_id| The unique identifier for the item the reporting dimension is for. | bigint |
|key| A reporting dimension key (e.g. "WBS") | ﻿character varying(50) |
|value| A reporting dimension value (e.g. OP011.896WIE) | character varying(100) |

<a name="auditing-changes"></a>
## Auditing Changes [(up)](#toc)

All changes to Items, Links and Item Types are recorded (via database triggers) on dedicated audit tables.

Audit tables contain the same attributes than the table they are auditing but add the following attributes:

|Name | Description | Data Type|
|---|---|---|
|operation | "I" for insert, "U" for update and "D" for delete. | character(1) |
|stamp| The time of the change. | timestamp |
|uderid| The user who performed the change. | text |
