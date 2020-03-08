# ox_item Resource

Creates, updates or destroys a configuration item.

## Example Usage

```hcl
resource "ox_item" "Item_1" {
  key         = "item_1"
  name        = "Item 1"
  description = "Item 1 Description"
  type        = "item_type_1"
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
  partition = "PARTITION_B"
}
```

## Argument Reference

* `key` - (Required) The item natural key (string).
* `name` - (Required) The item name (string).
* `description` - (Required) The item description (string).
* `type` - (Required) The item type natural key (string).
* `status` - (Optional) The item status (number).
* `meta` - (Optional) The item free JSON field.
* `txt` - (Optional) The item free text field.
* `tag` - (Optional) The item list of tags for searching (list of string).
* `attribute` - (Optional) The item key/value dictionary.
* `partition` - (Optional) The item logical partition (string).

<!-- ## Attribute Reference

* `attribute_name` - List attributes that this resource exports. -->
