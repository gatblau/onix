# ox_item_type_attribute Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration item type attribute data to be fetched for use elsewhere in Terraform configuration.

More information about item types can be found in the [Item Type Attribute Resource](../resources/ox_item_type_attr.md) section.

## Example Usage

```hcl
data "ox_item_type_attr" "attr_aws_instance_ram_data" {
  key = "ATTRIBUTE_AWS_INSTANCE_RAM"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the attribute.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| string | *The key name for the attribute as used in the attribute dictionary.* |
| `description`| string | *A meaningful description for the attribute.* |
| `item_type_key`| string | *The natural key that uniquely identifies the item type the attribute is for.* |
| `type` | string | *The attribute type, a free format string that should be understood by a client to decide how to validate the attribute value.* |
| `def_value`| string | *A free format string containing the default value for the attribute, if any.* |
| `required` | boolean | *A flag indicating whether the attribute is required.* |
| `regex`| string | *A regular expression used by a client to validate the attribute value.* |
| `version` | integer | *The version number of the attribute for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the attribute was created for the first time.* |
| `updated` | date & time | *The date and time the attribute was last updated.* |
| `changed_by` | string | *The user and role that last modified the attribute.* |
