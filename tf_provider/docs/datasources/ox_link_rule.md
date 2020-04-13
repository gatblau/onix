# ox_link_rule Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow configuration link rule data to be fetched for use elsewhere in Terraform configuration.

More information about link rules can be found in the [Link Rule Resource](../resources/ox_link_rule.md) section.

## Example Usage

```hcl
data "ox_link_rule" "aws_vpc_instance_rule_data" {
  key = "AWS_VPC->AWS_INSTANCE"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the link type.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| string | *The display name for the link rule.* |
| `description`| string | *A meaningful description for the link rule.* |
| `link_type_key` | string | *The natural key uniquely identifying the link type to which this rule applies.* |
| `start_item_type_key` | string | *The natural key uniquely identifying the item type of the starting item being linked.* |
| `end_item_type_key` | string | *The natural key uniquely identifying the item type of the ending item being linked.* |
| `version` | optional | integer | *The version number of the link type for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the link rule was created for the first time.* |
| `updated` | date & time | *The date and time the link rule was last updated.* |
| `changed_by` | string | *The user and role that last modified the link rule.* |