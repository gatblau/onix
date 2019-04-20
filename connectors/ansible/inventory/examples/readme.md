# Inventory Plugin Example <img src="../../../../docs/pics/ox.png" width="160" height="160" align="right">

This example shows how to use the inventory plugin to read inventories from the CMDB.

In order to use the plugin, ensure the plugin is installed in the control node as described [here](../readme.md).

Then follow the steps below.

### Install the Onix WAPI

The Onix WAPI is the Web API allowing reading and writing CI information to/from the database.
For the purpose of this example, it can be installed easily in containerised form.
Ensure you have [docker compose](https://docs.docker.com/compose/) installed in your machine and understand how to use it.

Run docker compose with this [docker-compose.yml](./../../../../install/container/docker-compose.yml) file.

Then run the following command from a linux terminal:

```bash
$ docker-compose up -d
```

### Import the inventory into the CMDB

In order to import Ansible inventory data into the CMDB, Swagger can be used to put an inventory [JSON payload](./inventory.json) 
on the */data* endpoint of the web service as follows:

1. Open the [Swagger UI](http://localhost:8080/swagger-ui.html#/web-api/createOrUpdateItemTreeUsingPUT)
2. Paste the payload [here](inventory.json) in the payload box
3. Execute the request to create the inventory in the CMDB
4. You should see the response of the web service showing no errors

### Create a tag for the inventory

The imported inventory needs to be tagged before it can be used.
The tag ensures that a version is fixed before it is used for deployment.

To take a tag:
1. Open the [Swagger UI](http://localhost:8080/swagger-ui.html#/web-api/createTagUsingPOST)
2. Paste the payload [here](tag.json) in the payload box
3. Put the request to create the tag in the CMDB

### Test the plugin

From this folder, type the following:
```bash
$ ansible-inventory -i onix_inventory.yml --graph
```
**NOTE**: the [onix_inventory.yml](./onix_inventory.yml) contains the configuration required by the plugin to connect to 
the CMDB.

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

[back to index](../modules/readme.md)