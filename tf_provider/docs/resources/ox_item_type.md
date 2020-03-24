# ox_item_type Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_item_type resource, creates, updates or destroys an item type definition associated to a model.

An item type defines the characteristics of a configuration item, such as the attributes that can record, the structure of JSON metadata in the item, etc.

An item type has to be defined before an item can be created.

## Example usage

In the example below, three item types are created to represent VPCs, Virtual Servers and a Storage Volumes:

```hcl
resource "ox_item_type" "AWS_VPC" {
  key         = "AWS_VPC"
  name        = "Virtual Private Cloud"
  description = "A logically isolated section of the AWS Cloud where AWS resources can be launched in a virtual network."
  model_key   = "AWS_EC2"
  managed     = false
}

resource "ox_item_type" "AWS_INSTANCE" {
  key         = "AWS_INSTANCE"
  name        = "AWS Instance"
  description = "A virtual server in the AWS Cloud."
  model_key   = "AWS_EC2"
  managed     = false
}

resource "ox_item_type" "AWS_EBS_VOLUME" {
  key         = "AWS_EBS_VOLUME"
  name        = "AWS Block Storage Volume"
  description = "A durable, block-level storage device that can be attached to one or more instances."
  model_key   = "AWS_EC2"
  managed     = false
}
```

## Argument reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the item type.* |
| `name`| required | string | *The display name for the item type.* |
| `description`| required | string | *A meaningful description for the item type.* |
| `filter`| optional | JSON | *Defines one or more filters, namely [JSON Path](https://goessner.net/articles/JsonPath/) expressions that allow the Web API to extract parts of the JSON metadata stored in a configuration item.* |
| `meta_schema` | optional | JSON | *The [JSON Schema](https://json-schema.org/) used to validate the JSON metadata stored in a configuration item.* |
| `model_key` | required | string | *The natural key uniquely identofying the model this item type is part of.* |

## Key dependencies



## Related entities

- [ox_link_type](ox_link_type.md)
- [ox_link_rule](ox_link_rule.md)
