# ox_model Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow data model data to be fetched for use elsewhere in Terraform configuration.

More information about configuration items can be found in the [Model Resource](../resources/ox_model.md) section.

## Example Usage

```hcl
data "ox_model" "EC2_data" {
  key = "EC2"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the model.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| string | *The display name for the model.* |
| `description`| string | *A meaningful description for the model.* |
| `partition` | string | *The natural key that identifies the logical partition the model is in.* |
| `version` | integer | *The version number for the model.* |
| `created` | date & time | *The date and time the model definition was first created.* |
| `updated` | date & time | *The date and time the model was last updated.* |
| `changed_by` | string | *The user and role that last modified the item.* |