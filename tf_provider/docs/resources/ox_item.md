# ox_item Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Creates, updates or destroys a configuration item.

A configuration item is a piece of data, of a particular type, which can be linked to other items of the same or different type. They carry versions and all their changes are recorded in the change history forming the basis of a configuration audit.

Items must be uniquely identified using its natural key so that they can be distinguished from all other configuration items.

Configuration Items can store data in either plain text, JSON or key/value pair(s) format.

JSON data can be validated using a JSON schema defined in the item type.

Equally, key/value pairs are validated through the attributes defined by the item type.

When storing sensitive information, a configuration item's text and JSON data can be automatically encrypted (at rest). The item type defines whether or not to encrypt text and/or JSON data.

## Example Usage

```hcl
resource "ox_item" "Item_1" {
  key         = "AWS_INSTANCE_ABC_001"
  name        = "AWS Instace - ABC-001"
  description = "AWS Instace - ABC-001 Description"
  type        = "AWS_INSTANCE"
  meta = {
     "hostvars": {
       "ansible_become": "yes",
       "openshift_master_overwrite_named_certificates": "True",
       "openshift_node_kubelet_args": {
         "image-gc-high-threshold": [
           "90"
         ],
         "image-gc-low-threshold": [
           "80"
         ]
       },
       "openshift_prometheus_pvc_size": "10Gi",
       "openshift_metrics_cassandra_pvc_storage_class_name": "glusterfs-storage-block",
       "openshift_image_tag": "v3.9.30"
  }
  txt = "Free format text here."
  attribute = {
    "RAM" : "3",
    "CPU" : "1"
  }
  tag = [ "tag1", "tag2", "tag3"]
  status = 3
  partition = "LOGISTICS_DEPT_PARTITION"
}
```

## Argument Reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the item.* |
| `name`| required | string | *The display name for the item.* |
| `description`| required | string | *A meaningful description for the item.* |
| `type` | required | string | *The natural key that uniquely identifies the [type of item](ox_item_type.md).* |
| `status` | optional | integer | *A number which describes the status the item is in. The default value is 0.* |
| `meta` | optional | json | *Stores any information in JSON format. It can be automatically encrypted if required.* |
| `txt` | optional | text | *Stores any information in text format. It can be automatically encrypted if required.* |
| `attribute` | optional | map of strings | *Stores zero or more key-value pairs that are defined in the item type.* |
| `tag` | optional | array of string | *Stores zero or more tags that can be used to classify or search for the item.* |
| `partition` | optional | string | *The natural key that identifies the logical partition the item is in. If no value is specified, the item is placed in the default instance partition (INS).* |
| `version` | optional | integer | *The version number of the item. If specified, optimistic locking is enabled: if the specified version is different than the stored version, no changes are made and a locking situation is assumed.* |

## Key dependencies

An item requires an [Item Type](ox_item_type.md) and [Partition](ox_partition.md) definitions.
If no [Partition](ox_partition.md) is specified, the Item is placed in the default instance partition (INS) by default. 

![Item](../pics/item.png)

## Related entities

- Item **is in** [Partition](ox_partition.md)
- Item **is of** [Item Type](ox_item_type.md)
- [Link](ox_link.md) **connects** Items

## Web API endpoints

This resource uses the following Web API endpoint: 

```bash
/item/{item_key}
```

The table below shows what methods are mapped to what operations in the terraform resource:

| **Method** | **Operation** |
|:---:|:---:|
| PUT | Create |
| GET | Read |
| PUT | Update |
| DELETE | Delete  |

