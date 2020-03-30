# ox_item_type Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration item type data to be fetched for use elsewhere in Terraform configuration.

More information about item types can be found in the [Item Type Resource](../resources/ox_item_type.md) section.

## Example Usage

```hcl
data "ox_item_type" "aws_instance_data" {
  key = "AWS_INSTANCE"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the item type.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| string | *The display name for the item type.* |
| `description` | string | *A meaningful description for the item type.* |
| `model_key` | string | *The natural key uniquely identifying the model this item type is part of.* |
| `filter` | JSON | *Defines one or more filters, namely [JSON Path](https://goessner.net/articles/JsonPath/) expressions that allow the Web API to extract parts of the JSON metadata stored in a configuration item. The format of the filter is described in the notes section below.* |
| `meta_schema` | JSON | *The [JSON Schema](https://json-schema.org/) used to validate the JSON metadata stored in a configuration item's meta attribute.* |
| `notify_change` | boolean | *Determines whether notification events should be sent by the Web API when items of this type are created, updated or deleted. The default value is false.* |
| `tag` | string array | *A list of tags used for searching and classifying the item type.* |
| `encrypt_meta` | boolean | *A flag indicating whether the meta attribute of the configuration item of this type should have encryption of data at rest.* |
| `encrypt_txt` | boolean | *A flag indicating whether the txt attribute of the configuration item of this type should have encryption of data at rest.* |
| `managed` | boolean | *A flag indicating whether the item type is managed by a third party process. The default value is false, indicating the type can be updated by the user interface or Terraform provider clients.* |
| `version` | integer | *The version number of the item type for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the item type was created for the first time.* |
| `updated` | date & time | *The date and time the item type was last updated.* |
| `changed_by` | string | *The user and role that last modified the item type.* |