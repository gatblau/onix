# ox_item Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration item data to be fetched for use elsewhere in Terraform configuration. 

Use of the *ox_item* data source allows a Terraform configuration to make use of information defined in Onix.

## Example Usage

```hcl
data "ox_item" "item_1_data" {
  key = "item_1_key"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the item.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `id` | string | *The surrogate key that uniquely identifies the item.* |
| `key` | string | *The natural key that uniquely identifies the item.* |
| `name`| string | *The display name for the item.* |
| `description`| string | *A meaningful description for the item.* |
| `type` | string | *The natural key that uniquely identifies the type of item.* |
| `status` | integer | *A number which describes the status the item is in. The default value is 0.* |
| `meta` | json | *Stores any information in JSON format. It can be automatically encrypted if required.* |
| `txt` | text | *Stores any information in text format. It can be automatically encrypted if required.* |
| `attribute` | map[string, object] | *Stores zero or more key-value pairs that are defined in the item type.* |
| `tag` | array of string | *Stores zero or more tags that can be used to classify or search for the item.* |
| `partition` | string | *The natural key that identifies the logical partition the item is in. If no value is specified, the item is placed in the default instance partition (INS).* |
| `encKeyIx` | integer | *The index of the encryption key used to encrypt text and JSON data in the item. It can be 0 if no encryption key was used, 1 or 2 if encryption was used.* |
| `version` | integer | *The version number for the item. Every time a change is made to an item, the version number is automatically incremented. The version number is used to enable optimistic concurrency locking.* |
| `created` | date & time | *The date and time the item was created for the first time.* |
| `updated` | date & time | *The date and time the item was last updated.* |
