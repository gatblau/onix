# Onix Inventory Plugin for Ansible

Onix provides an [inventory plugin](./onix.py) which allows users to point at a Onix CMDB to compile the inventory of 
hosts that Ansible uses to target tasks.

This plugin works in a similar way to the Ansible Tower inventory plugin.

In order to use the plugin, make sure the following apply:

## Plugin location

The plugin is located in the default inventory plugin path in the control machine, typically under:
 - ~/.ansible/plugins/inventory' or
 - /usr/share/ansible/plugins/inventory'
 
 In order to determine the value of the default path in your system do the following:
 
```bash
$ ansible-config dump | grep DEFAULT_INVENTORY_PLUGIN_PATH
```

## Plugin configuration information

The plugin reads configuration information from either of two sources:

1. **YAML file**: called either onix.yml or onix_inventory.yml and loaded using ansible-playbook -i option passing the file path.
2. **Environment variables**: if the option passed to ansible-playbook -i is **@onix_inventory**
3. **A combination of the above**: if some of the variables are missing from the onix_inventory.yml file the plugin will attempt
to read them from environment variables.

The following table describes the configuration variables required:

|Name | Description | Env Variable | Example |
|---|---|---|---|
|**plugin**| The name of the plpugin. Set to onix.| N/A| onix
|**host**| The network address of the Onix WAPI service. | OX_HOST | localhost:8080 |
|**username**| The user to authenticate with the Onix WAPI. | OX_USERNAME| admin |
|**password**| The password to authenticate with the Onix WAPI. | OX_PASSWORD | 0n1x |
|**inventory_key**| The natural key for the inventory in the CMDB. | OX_INVENTORY_KEY | inventory_01 |
|**inventory_version**| The tag associated to a version of the inventory in the CMDB. | OX_INVENTORY_TAG | v1 |
|**verify_ssl**| Whether to verify TLS keys when calling the Onix WAPI service. | OX_VERIFY_SSL | false |

## Enabling the plugin

Most inventory plugins shipped with Ansible are disabled by default and need to be whitelisted in your ansible.cfg file 
in order to function. This is how the whitelist should look like in the config file to enable the onix plugin:

In the **ansible.cfg** file:
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

Running:
```bash
$ ansible-inventory -i onix_inventory.yml --list -vvv
```
should display the inventory structure loaded from Onix.