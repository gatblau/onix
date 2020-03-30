# ox_link_rule Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_link_rule resource, creates, updates or destroys rules that dictate how [items](ox_item.md) should be connected to [links](ox_link.md).

Link rules must be defined to allow [items](ox_item.md) of a particular [item type](ox_item_type.md) to connect using [links](ox_link.md) of a particular [link type](ox_link_type.md).

## Example usage

In the example below, two link rules are created to allow links between:

- VPCs and Virtual Servers and
- Virtual Servers and Storage Volumes

```hcl
resource "ox_link_rule" "AWS_EC2_VPC_INSTANCE" {
  key                 = "AWS_VPC->AWS_INSTANCE"
  name                = "VPCs To VMs Link Rule"
  description         = "Allows to link AWS VPCs to Instances."
  link_type_key       = "AWS_EC2_LINK"
  start_item_type_key = "AWS_VPC"
  end_item_type_key   = "AWS_INSTANCE"
}

resource "ox_link_rule" "AWS_EC2_INSTANCE_EBS_VOLUME" {
  key                 = "AWS_INSTANCE->AWS_EBS_VOLUME"
  name                = "Instances To EBS Volumes Link Rule"
  description         = "Allows to link AWS Instances to EBS Volumes."
  link_type_key       = "AWS_EC2_LINK"
  start_item_type_key = "AWS_INSTANCE"
  end_item_type_key   = "AWS_EBS_VOLUME"
}
```

## Argument reference

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | *required* | string | *The natural key that uniquely identifies the link rule.* |
| `name`| *required* | string | *The display name for the link rule.* |
| `description`| *required* | string | *A meaningful description for the link rule.* |
| `link_type_key` | *required* | string | *The natural key uniquely identifying the link type to which this rule applies.* |
| `start_item_type_key` | *required* | string | *The natural key uniquely identifying the item type of the starting item being linked.* |
| `end_item_type_key` | *required* | string | *The natural key uniquely identifying the item type of the ending item being linked.* |
| `version` | optional | integer | *The version number of the link type for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |

## Key dependencies

A link rule requires a [link type](ox_link_type.md), and one or two [item types](ox_item_type.md).

![Link Rule](../pics/link_rule.png)

## Related entities

- Link Rule **allows links between items of** [Item Type](ox_item_type.md)
- Link Rule **is for** [Link Type](ox_item_type.md)

## Web API endpoints

This resource uses the following Web API endpoint: 

```bash
/linkrule/{link_rule_key}
```

The table below shows what methods are mapped to what operations in the terraform resource:

| **Method** | **Operation** |
|:---:|:---:|
| PUT | Create |
| GET | Read |
| PUT | Update |
| DELETE | Delete  |