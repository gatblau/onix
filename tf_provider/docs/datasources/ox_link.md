# ox_link Data Source <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration item link data to be fetched for use elsewhere in Terraform configuration.

More information about configuration items can be found in the [Link Resource](../resources/ox_link.md) section.

## Example Usage

```hcl
data "ox_link" "vpc01_vm01_link_data" {
  key = "vpc01_vm01_link"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the link between two items.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `description`| string | *A meaningful description for the link.* |
| `type` | string | *The natural key that uniquely identifies the type of link.* |
| `meta` | json | *Stores any information in JSON format. It can be automatically encrypted if required.* |
| `txt` | text | *Stores any information in text format. It can be automatically encrypted if required.* |
| `attribute` | map[string, object] | *Stores zero or more key-value pairs that are defined in the [link type](../resources/ox_link.md).* |
| `tag` | array of string | *Stores zero or more tags that can be used to classify or search for the link.* |
| `version` | integer | *The version number for the link. Every time a change is made to a link, its version number is automatically incremented. The version number is used to enable optimistic concurrency locking.* |
| `created` | date & time | *The date and time the link was first created.* |
| `updated` | date & time | *The date and time the link was last updated.* |
| `changed_by` | string | *The user and role that last modified the item.* |
