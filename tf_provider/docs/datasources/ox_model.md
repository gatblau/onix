# ox_model Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow data model data to be fetched for use elsewhere in Terraform configuration. 

## Example Usage

```hcl
data "ox_model" "model_1_data" {
  key = "model_1_key"
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
| `id` | string | *The surrogate key that uniquely identifies the model.* |
| `key` | string | *The natural key that uniquely identifies the model.* |
| `name`| string | *The display name for the model.* |
| `description`| string | *A meaningful description for the model.* |
| `partition` | string | *The natural key that identifies the logical partition the model is in.* |
| `version` | integer | *The version number for the model.* |
| `created` | date & time | *The date and time the model definition was first created.* |
| `updated` | date & time | *The date and time the model was last updated.* |
