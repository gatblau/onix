# ox_link Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Creates, updates or destroys a link between two configuration items.

## Example Usage

```hcl
resource "ox_link" "Link_1" {
  key            = "link_1"
  description    = "link 1 description"
  type           = "Test_Link_Type"
  start_item_key = "Item_1"
  end_item_key   = "Item_2"
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
    "TEAM" : "Blue",
    "CATEGORY" : "Social"
  }
  tag = [ "tag1", "tag2", "tag3"]
}
```

## Argument Reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the link.* |
| `description`| required | string | *A meaningful description for the link.* |
| `type` | required | string | *The natural key that uniquely identifies the [type of link](ox_link_type.md).* |
| `meta` | optional | json | *Stores any information in JSON format. It can be automatically encrypted if required.* |
| `txt` | optional | text | *Stores any information in text format. It can be automatically encrypted if required.* |
| `attribute` | optional | map of strings | *Stores zero or more key-value pairs that are defined in the item type.* |
| `tag` | optional | array of string | *Stores zero or more tags that can be used to classify or search for the item.* |
| `version` | optional | integer | *The version number of the item. If specified, optimistic locking is enabled: if the specified version is different than the stored version, no changes are made and a locking situation is assumed.* |
