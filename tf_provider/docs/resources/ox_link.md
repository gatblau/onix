# ox_link Resource <img src="../../docs/pics/ox.png" width="200" height="200" align="right">

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

* `key` - (Required) The link natural key (string).
* `description` - (Required) The link description (string).
* `type` - (Required) The link type natural key (string).
* `meta` - (Optional) The item free JSON field.
* `txt` - (Optional) The item free text field.
* `tag` - (Optional) The item list of tags for searching (list of string).
* `attribute` - (Optional) The item key/value dictionary.

<!-- ## Attribute Reference

* `attribute_name` - List attributes that this resource exports. -->
