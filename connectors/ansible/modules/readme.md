# Ansible Modules <img src="../../../docs/pics/ox.png" width="125" height="125" align="right">

In order to maintain accurate information in the CMDB it is important that:
- Changes are recorded in the CMDB as soon as configuration changes are made (If changes are driven by Ansible, then it is important to update the CMDB automatically whilst Ansible performs any configuration changes). Or
- Changes are driven by information in the CMDB. Or
- Both of the above cases.

In order to facilitate this, a set of Ansible modules are provided as follows:

- [Ansible Modules](#ansible-modules)
  - [ox_setup](#oxsetup)
  - [ox_item_type](#oxitemtype)
  - [ox_link_type](#oxlinktype)
  - [ox_link_rule](#oxlinkrule)
  - [ox_item](#oxitem)
  - [ox_link](#oxlink)
  - [ox_import](#oximport)
  - [ox_query](#oxquery)

For particular examples of how to use these modules see [here](./examples/readme.md).

### State: Present or Absent

Modules have a "state" property which always default to "present".

If **state: "absent"** is specified, then the item, item type or link will be removed from the database.
 

### Idempotency

Modules are idempotent, which means if the module that has created an item/item type/link is run for a second time, changes will not be made to the database.

If however, any of the item/item type/link properties has changed, tan uodated version will be written to the database.

 
### How to use the Ansible modules
 
For an example of how to use the above modules, take a look at the [playbook here](../ansible/site.yml).
 
To execute the playbook run the following command from the [ansible](../ansible) folder:

```bash
$ ansible-playbook -i inventory site.yml -vvv
```

**NOTE**: it is assumed Onix Service, Onix Database and Keycloak are running in the localhost under the default ports.
If this is not the case, update the variables in the [inventory](../ansible/inventory) file accordingly.

### How to install the Ansible modules

#### In the same location of the playbooks project

Copy the files [here](../../ansible/modules/) to "binary" folder under the folder where the playbook using the modules is. 

#### In a shared location

To share the modules across multiple projects, add an entry to the **/etc/ansible/ansible.cfg** file pointing to a shared library location as follows:

```bash
library = /usr/share/ansible/library
```

## Modules

<a name="login"></a>
### [ox_setup](../modules/ox_setup.py)

The **ox_setup** module is used to setup the location and authentication information to 
connect to the Onix WAPI.

It has to be called **before** calling any other modules.

Use it as follows:
```yaml
- name: configure access to Onix
  ox_setup:
    uri: "http://localhost:8080" 
    username: "admin" 
    password: "0n1x" 
    auth_mode: "basic"
```
**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **username** | the username for basic authentication | yes |
| **password** | the password for basic authentication | yes |
| **auth_mode** | the authentication mode to use (i.e. none, basic, openid) | no (default to none) |

**NOTE:** the module creates two facts meant to be used as input variables on the other modules as follows:

| name | description |
|---|---|
| **ox_uri** | the URI of the Onix WAPI |
| **ox_token** | the token used to authenticate with the Onix WAPI |


<a name="item_type"></a>
### [ox_item_type](../modules/ox_item_type.py)

The **ox_item_type** module is used to create/update or delete item types in the CMDB.

Use it as follows:
```yaml
- name: Creates the Application Item Type
  ox_item_type:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    key: "APPLICATION"
    name: "Software Application"
    description: "A Software Application."
```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **key** | the natural key for the item type | yes |
| **name** | a user friendly name for the item type | no |
| **description** | the item type description | no |
| **state** | *'present'* to create/update the item type; or *'absent'* to delete the item type. | no (default to *present*) |

<a name="link_type"></a>
### [ox_link_type](../modules/ox_link_type.py)

The **ox_link_type** module is used to create/update or delete link types in the CMDB.

Use it as follows:
```yaml
- name: Creates the Application Link Type
  ox_link_type:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    key: "APPLICATION"
    name: "Software Application Link."
    description: "Links application services."
```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **key** | the natural key for the link type | yes |
| **name** | a user friendly name for the link type | no |
| **description** | the link type description | no |
| **state** | *'present'* to create/update the link type; or *'absent'* to delete the link type. | no (default to *present*) |

<a name="link_rule"></a>
### [ox_link_rule](../modules/ox_link_rule.py)

The **ox_link_rule** module is used to create/update or delete link rules in the CMDB.

Use it as follows:
```yaml
- name: Creates the Application Runtime to Host Link Rule
  ox_link_rule:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    key: "APPLICATION-RUNTIME->HOST"
    linkTypeKey: "APPLICATION"
    startItemTypeKey: "APPLICATION-RUNTIME"
    endItemTypeKey: "HOST"
    name: "Software Application Runtime to Host Link Rule."
    description: "Allows linking application runtime items with host items."

```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **key** | the natural key for the link rule | yes |
| **linkTypeKey** | the natural key for the link type the rule is for | yes |
| **startItemTypeKey** | the natural key for the type of the item that should start the link. | no |
| **endItemTypeKey** | the natural key for the type of the item that should end the link. | no |
| **name** | a user friendly name for the link type | no |
| **description** | the link type description | no |
| **link** | the link type description | no |
| **state** | *'present'* to create/update the link type; or *'absent'* to delete the link type. | no (default to *present*) |

<a name="item"></a>
### [ox_item](../modules/ox_item.py)

The **ox_item** module is used to create/update or delete configuration items in the CMDB.

Use it as follows:
```yaml
- name: Creates a configuration for Onix application Data Service
  ox_item:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    key: "DATA-SERVICE-ONIX-DB"
    name: "ONIX Data Service"
    description: "Onix Data Service"
    type: "DATA-SERVICE"
    meta: {
    }
    status: 1
    tag:
    - "onix"
    - "db"
    attrValid:
      WBS: "EU-00023.100002.984"
      PROJECT: "TheOnixProject"
```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **key** | the natural key for the item | yes |
| **name** | a user friendly name for the item | no |
| **description** | the item description | no |
| **type** | the item type key created by the ox_item_type module | no |
| **meta** | an arbitrary json object associated with the item | no |
| **status** | an arbitrary flag indicating the status of the item | no |
| **tag** | a list of tags for searching | no |
| **attrValid** | a map of key-value pairs for searching | no |
| **state** | *'present'* to create/update the item type; or *'absent'* to delete the item type. | no (default to *present*) |

<a name="link"></a>
### [ox_link](../modules/ox_link.py)

The **ox_link** module is used to create/update or delete links between existing configuration items in the CMDB.

Use it as follows:
```yaml
- name: Creates a link between Spring Boot and Host B
  ox_link:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    key: "RUNTIME-SPRING-BOOT-HOSTB"
    description: "Spring Boot is deployed on Host B."
    type: "APPLICATION"
    startItemKey: "RUNTIME-SPRING-BOOT"
    endItemKey: "HOST-B"
    tag:
    - "runtime"
    attrValid:
      WBS: "EU-00023.100002.984"
      PROJECT: "TheOnixProject"
    meta: {
      runtime: "Spring Boot",
      version: "1.5.10.RELEASE",
    }

```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **key** | the natural key for the link | yes |
| **name** | a user friendly name for the link | no |
| **description** | the link description | no |
| **type** | the link type key created by the ox_link_type module | no |
| **startItemKey** | the key of the item that starts the link | no |
| **endItemKey** | the key of the item that ends the link | no |
| **meta** | an arbitrary json object associated with the item | no |
| **tag** | a list of tags for searching | no |
| **attrValid** | a map of key-value pairs for searching | no |
| **state** | *'present'* to create/update the item type; or *'absent'* to delete the item type. | no (default to *present*) |

<a name="import"></a>
### [ox_import](../modules/ox_import.py)

This module imports the configuration data in a json file.

It is convenient when a set of configuration data has to be imported at once, particularly when such data can be 
templated using for example the [Ansible template](https://docs.ansible.com/ansible/latest/modules/template_module.html) 
module before posting it to the CMDB.

Use it as follows:
```yaml
- name: import configuration
  ox_import:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    src: "path_to_configuration_file.json"
```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **src** | the path (relative or absolute) to a json file containing the configuration data to be imported. | yes |

Example configuration file:
```json
{
    "models": [
    {
        "key": "APP_MODEL",
        "name": "Application Meta Model.",
        "description": "Describes the item and link types required to represent an Application in the CMDB."
    }],
    "itemTypes": [
    {
        "key": "APPLICATION",
        "name": "Software Application",
        "description": "A Software Application.",
        "modelKey": "APP_MODEL"
    },
    {
        "key": "WEB_SERVICE",
        "name": "Web Service",
        "description": "A web service that is part of an application.",
        "modelKey": "APP_MODEL"
    }],
    "linkTypes": [
    {
        "key": "APPLICATION",
        "name": "Software Application Link",
        "description": "Links application services.",
        "modelKey": "APP_MODEL"
    }],
    "linkRules": [
    {
        "key": "APPLICATION->WEB-SERVICE",
        "linkTypeKey": "APPLICATION",
        "startItemTypeKey": "APPLICATION",
        "endItemTypeKey": "WEB_SERVICE",
        "name": "Software Application to Web Service Link Rule.",
        "description": "Allows linking application items with web-service items."
    }]
}
```
<a name="query"></a>
### [ox_query](../modules/ox_query.py)

This module retrieves configuration data from the CMDB.

It uses the natural key to retrieve information about one of the following objects:
 - item 
 - link
 - item type
 - link type
 - link rule
 - model
 
Use it as follows:
```yaml
- name: query item data
  ox_query:
    uri: "{{ ox_uri }}"
    token: "{{ ox_token }}"
    type: "item" # can also use link, link_type, item_type, model, link_rule
    key: "NODE_1"
  register: result
```
**NOTE**: *ox_uri* and *ox_token* are facts produced by the *ox_setup* module.

**where:**

| name | description | required |
|---|---|---|
| **uri** | the URL of the Onix CMDB. | yes |
| **token** | the access token for the service. | yes |
| **type** | the type of the object to be queried. Allowed types are "item", "link", "item_type", "link_type", "link_rule" and "model" | yes |
| **key** | the natural key of the object to query | yes |

Once the result has been registered, it can be output as follows:

```yaml
- name: print whole query result
  debug:
    var: result

- name: print item metadata
  debug:
    var: result.meta

- name: print item attributes
  debug:
    var: result.attrValid

- name: print item tags
  debug:
    var: result.tag
```
([up](../readme.md))