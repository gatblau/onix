# ox_item Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration item data to be fetched for use elsewhere in Terraform configuration.

More information about configuration items can be found in the [Item Resource](../resources/ox_item.md) section.

## Example Usage

```hcl
data "ox_item" "aws_instance_abc_001_data" {
  key = "AWS_INSTANCE_ABC_001"
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
| `changed_by` | string | *The user and role that last modified the item.* |
