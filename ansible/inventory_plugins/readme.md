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

## Trying it out

In order to try the plugin out, do the following:

### Install the Onix WAPI

Install onix wapi by running [docker compose](https://docs.docker.com/compose/) with the configuration 
file located [here](./../../install/container/docker-compose.yml).

```bash
$ docker-compose up -d
```

### Import the inventory into the CMDB

1. Open the [Swagger UI](http://localhost:8080/swagger-ui.html#/web-api/createOrUpdateItemTreeUsingPUT)
2. Paste the payload [here](inventory.json) in the payload box
3. Put the request to create the inventory in the CMDB

### Create a snapshot for the inventory

The imported inventory needs to be snapshoted before it can be used.
The snapshot ensures that a version is fixed before it is used for deployment.

To take a snapshot:
1. Open the [Swagger UI](http://localhost:8080/swagger-ui.html#/web-api/createSnapshotUsingPOST)
2. Paste the payload [here](snapshot.json) in the payload box
3. Put the request to create the snapshot in the CMDB

### Test the plugin

From this folder, type the following:
```bash
$ ansible-inventory -i onix_inventory.yml --graph
```
The output should look like:
```text
@all:
  |--@test_inventory::OSEv3:
  |  |--@test_inventory::infra:
  |  |  |--test_inventory::infra1
  |  |  |--test_inventory::infra2
  |  |--@test_inventory::masters:
  |  |  |--test_inventory::master1
  |  |  |--test_inventory::master2
  |--@test_inventory::compute:
  |  |--test_inventory::compute1
  |  |--test_inventory::compute2
  |--@ungrouped:
```
Now type the following command to get json with all the inventory data as loaded in Ansible:
```bash
ansible-inventory -i onix_inventory.yml --list
```
the following output should be displayed in the terminal:
```json
{
    "_meta": {
        "hostvars": {
            "test_inventory::compute1": {
                " ansible_host": "host0003.example.com", 
                " openshift_node_labels": "\"{'node': 'true', 'region': 'primary', 'zone': 'az2', 'site': 'b'}\""
            }, 
            "test_inventory::compute2": {
                " ansible_host": "host0005.example.com", 
                " openshift_node_labels": "\"{'node': 'true', 'region': 'primary', 'zone': 'az3', 'site': 'c'}\""
            }, 
            "test_inventory::infra1": {
                " ansible_host": "host0011.example.com", 
                " openshift_node_labels": "\"{'region': 'infra', 'zone': 'az2', 'site': 'b'}\"", 
                "ansible_become": "yes", 
                "openshift_image_tag": "v3.9.30", 
                "openshift_master_overwrite_named_certificates": "True", 
                "openshift_metrics_cassandra_pvc_storage_class_name": "glusterfs-storage-block", 
                "openshift_node_kubelet_args": {
                    "image-gc-high-threshold": [
                        "90"
                    ], 
                    "image-gc-low-threshold": [
                        "80"
                    ]
                }, 
                "openshift_prometheus_pvc_size": "10Gi"
            }, 
            "test_inventory::infra2": {
                " ansible_host": "host0012.example.com", 
                " openshift_node_labels": "\"{'region': 'infra', 'zone': 'az3', 'site': 'c'}\"", 
                "ansible_become": "yes", 
                "openshift_image_tag": "v3.9.30", 
                "openshift_master_overwrite_named_certificates": "True", 
                "openshift_metrics_cassandra_pvc_storage_class_name": "glusterfs-storage-block", 
                "openshift_node_kubelet_args": {
                    "image-gc-high-threshold": [
                        "90"
                    ], 
                    "image-gc-low-threshold": [
                        "80"
                    ]
                }, 
                "openshift_prometheus_pvc_size": "10Gi"
            }, 
            "test_inventory::master1": {
                " ansible_host": "host0002.example.com", 
                " openshift_node_labels": "\"{'node': 'false', 'region': 'master', 'zone': 'az2'}\"", 
                "ansible_become": "yes", 
                "openshift_image_tag": "v3.9.30", 
                "openshift_master_overwrite_named_certificates": "True", 
                "openshift_metrics_cassandra_pvc_storage_class_name": "glusterfs-storage-block", 
                "openshift_node_kubelet_args": {
                    "image-gc-high-threshold": [
                        "90"
                    ], 
                    "image-gc-low-threshold": [
                        "80"
                    ]
                }, 
                "openshift_prometheus_pvc_size": "10Gi"
            }, 
            "test_inventory::master2": {
                " ansible_host": "host0004.example.com", 
                " openshift_node_labels": "\"{'node': 'false', 'region': 'master', 'zone': 'az3'}\"", 
                "ansible_become": "yes", 
                "openshift_image_tag": "v3.9.30", 
                "openshift_master_overwrite_named_certificates": "True", 
                "openshift_metrics_cassandra_pvc_storage_class_name": "glusterfs-storage-block", 
                "openshift_node_kubelet_args": {
                    "image-gc-high-threshold": [
                        "90"
                    ], 
                    "image-gc-low-threshold": [
                        "80"
                    ]
                }, 
                "openshift_prometheus_pvc_size": "10Gi"
            }
        }
    }, 
    "all": {
        "children": [
            "test_inventory::OSEv3", 
            "test_inventory::compute", 
            "ungrouped"
        ]
    }, 
    "test_inventory::OSEv3": {
        "children": [
            "test_inventory::infra", 
            "test_inventory::masters"
        ]
    }, 
    "test_inventory::compute": {
        "hosts": [
            "test_inventory::compute1", 
            "test_inventory::compute2"
        ]
    }, 
    "test_inventory::infra": {
        "hosts": [
            "test_inventory::infra1", 
            "test_inventory::infra2"
        ]
    }, 
    "test_inventory::masters": {
        "hosts": [
            "test_inventory::master1", 
            "test_inventory::master2"
        ]
    }, 
    "ungrouped": {}
}
```