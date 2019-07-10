# Inventory Plugin <img src="../../../docs/pics/ox.png" width="160" height="160" align="right">

The [inventory plugin](./onix.py) allows users to point at a Onix CMDB to retrieve an Ansible inventory. 

The plugin behaves in a similar way to the Ansible Tower inventory plugin.

In order to use the plugin, make sure the following apply:

## Plugin location

The plugin should be located in the default inventory plugin path in the control machine, typically under:
 - ~/.ansible/plugins/inventory' or
 - /usr/share/ansible/plugins/inventory'
 
 In order to determine the default path in your system do the following:
 
```bash
$ ansible-config dump | grep DEFAULT_INVENTORY_PLUGIN_PATH
```

## Plugin configuration information

The plugin reads configuration information from either of two sources:

1. **YAML file**: called either onix.yml or onix_inventory.yml and loaded using the *ansible-playbook -i* option passing the file path.
2. **Environment variables**: if the option passed to *ansible-playbook -i* option is **@onix_inventory**, then configuration is taken from environment variables.
3. **A combination of the above**: if some of the variables are missing from the onix_inventory.yml file the plugin will attempt
to read them from environment variables.

The following table describes the configuration variables required:

|Attribute | Description | Environment Var. | Example value|
|---|---|---|---|
|**plugin**| The name of the plugin. Set to onix.| N/A| onix |
|**host**| The network address of the Onix WAPI service. | OX_HOST | localhost:8080 |
|**username**| The user to authenticate with the Onix WAPI. | OX_USERNAME| admin |
|**password**| The password to authenticate with the Onix WAPI. | OX_PASSWORD | 0n1x |
|**inventory_key**| The natural key for the inventory in the CMDB. | OX_INVENTORY_KEY | inventory_01 |
|**inventory_version**| The tag associated to a version of the inventory in the CMDB. | OX_INVENTORY_TAG | v1 |
|**verify_ssl**| Whether to verify TLS keys when calling the Onix WAPI service. | OX_VERIFY_SSL | false |
|**auth_mode**| The approach used to authenticate with the Web API. Possible values are __none__, __basic__ (basic authentication) or __oidc__ (OpenId/OAuth 2.0). | OX_AUTH_MODE | basic |
|**client_id**| The unique identifier for the Onix Web API application as defined in the OAuth 2.0 server. It is only required if auth_mode is set to oidc. | OX_CLIENT_ID | dece7re....sxsxndj |
|**secret** | A secret known only to the application and the authorisation server. It is only required if auth_mode is set to oidc. | OX_SECRET | SXOUND...xssiuxnSIQ |
|**token_uri**| The OAuth 2.0 server endpoint where the ox provider exchanges the user credentials, client ID and client secret, for an access token. It is only required if auth_mode is set to oidc. | OX_TOKEN_URI | https://dev-1234.okta.com/oauth2/default/v1/token |

## Enabling the plugin

Most inventory plugins shipped with Ansible are disabled by default and need to be whitelisted in your ansible.cfg file 
in order to function. 

This is how the whitelist should look like in the [ansible.cfg](https://docs.ansible.com/ansible/latest/cli/ansible-config.html) 
config file to enable the onix plugin:
 
```bash
[inventory]
enable_plugins = onix, host_list, script, yaml, ini, auto
```

## Validating the plugin is enabled

Running: 
```bash
$  ansible-config dump | grep INVENTORY_ENABLED
```
should display a list of enabled plugings including onix.

## Deploying the Inventory meta model

In order for the inventory plugin to be able to retrieve inventory information, the Ansible Inventory meta model has to be imported into Onix.
The meta model contains all the item types, link types and link rules required to represent the inventory data in the CMDB.

To import the meta model follow the steps shown [here](../../../docs/models/readme.md).

## Trying it out

In order to try the plugin out, see an example in the [examples folder](./examples/readme.md).

### NOTE
This plugin is based on the Tower plugin by Matthew Jones and Yunfan Zhang as shown [here](https://github.com/ansible/ansible/blob/stable-2.7/lib/ansible/plugins/inventory/tower.py).

Information on how to use the Tower plugin can be found [here](https://docs.ansible.com/ansible/latest/plugins/inventory/tower.html).

([back to index](../../readme.md))