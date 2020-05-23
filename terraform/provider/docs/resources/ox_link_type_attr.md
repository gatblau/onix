# ox_link_type_attribute Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_link_type_attribute resource, creates, updates or destroys attributes of a link type.

An link type attribute defines the name and validation of attributes placed in a link *attribute* key-value dictionary.

An link type attribute requires an link type to exist first.

## Example usage

In the example below, two attributes are created to store CPU and RAM values for AWS Instances:

```hcl
resource "ox_link_type_attr" "aws_ec2_link_ram_attr" {
  key           = "ATTRIBUTE_EC2_LINK_WEIGHT"
  name          = "WEIGHT"
  description   = "The weight to calculate an aggregated rating for linked items."
  link_type_key = "AWS_EC2_LINK"
  type          = "decimal"
  def_value     = "0.2"
  managed       = false
}
```

## Argument reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | *required* | string | *The natural key that uniquely identifies the attribute.* |
| `name`| *required* | string | *The key name for the attribute as used in the link type attribute dictionary.* |
| `description`| *required* | string | *A meaningful description for the attribute.* |
| `link_type_key`| *required* | string | *The natural key that uniquely identifies the link type the attribute is for.* |
| `type` | *required* | string | *The attribute type, a free format string that should be understood by a client to decide how to validate the attribute value.* |
| `def_value`| optional | string | *A free format string containing the default value for the attribute, if any.* |
| `required` | optional | boolean | *A flag indicating whether the attribute is required.* |
| `regex`| optional | string | *A regular expression used by a client to validate the attribute value.* |
| `version` | optional | integer | *The version number of the attribute for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |

## Key dependencies

An item type belongs in a model and therefore, a model should exist first and be specified by the *model_key* attribute.

![Link Type Attribute](../pics/link_type_attr.png)

## Related entities

- [Link Type](ox_link_type.md) **has** Link Type Attribute(s)

## Web API endpoints

This resource uses the following Web API endpoint: 

```bash
/linktype/{link_type_key}/attribute/{attribute_key}
```

The table below shows what methods are mapped to what operations in the terraform resource:

| **Method** | **Operation** |
|:---:|:---:|
| PUT | Create |
| GET | Read |
| PUT | Update |
| DELETE | Delete  |