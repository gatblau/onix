# ox_item_type_attribute Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_item_type_attribute resource, creates, updates or destroys attributes of an item type.

An item type attribute defines the name and validation of attributes placed in a configuration item *attribute* key-value dictionary.

An item type attribute requires an item type to exist first.

## Example usage

In the example below, two attributes are created to store CPU and RAM values for AWS Instances:

```hcl
resource "ox_item_type_attr" "aws_instance_ram_attr" {
  key           = "ATTRIBUTE_AWS_INSTANCE_RAM"
  name          = "RAM"
  description   = "GB of RAM"
  item_type_key = "AWS_INSTANCE"
  type          = "integer"
  def_value     = "2"
  managed       = false
}

resource "ox_item_type_attr" "aws_instance_cpu_attr" {
  key           = "ATTRIBUTE_AWS_INSTANCE_CPU"
  name          = "CPU"
  description   = "No of CPU"
  item_type_key = "AWS_INSTANCE"
  type          = "integer"
  def_value     = "1"
  managed       = false
}
```

## Argument reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | *required* | string | *The natural key that uniquely identifies the attribute.* |
| `name`| *required* | string | *The key name for the attribute as used in the attribute dictionary.* |
| `description`| *required* | string | *A meaningful description for the attribute.* |
| `item_type_key`| *required* | string | *The natural key that uniquely identifies the item type the attribute is for.* |
| `type` | *required* | string | *The attribute type, a free format string that should be understood by a client to decide how to validate the attribute value.* |
| `def_value`| optional | string | *A free format string containing the default value for the attribute, if any.* |
| `required` | optional | boolean | *A flag indicating whether the attribute is required.* |
| `regex`| optional | string | *A regular expression used by a client to validate the attribute value.* |
| `managed` | optional | boolean | *A flag indicating whether the attribute is managed by a third party process. The default value is false, indicating the type can be updated by the user interface or Terraform provider clients.* |
| `version` | optional | integer | *The version number of the attribute for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |

## Key dependencies

An item type belongs in a model and therefore, a model should exist first and be specified by the *model_key* attribute.

![Item Type Attribute](../pics/item_type_attr.png)

## Related entities

- [Item Type](ox_item_type.md) **has** Item Type Attribute(s)

## Web API endpoints

This resource uses the following Web API endpoint: 

```bash
/itemtype/{link_type_key}/attribute/{attribute_key}
```

The table below shows what methods are mapped to what operations in the terraform resource:

| **Method** | **Operation** |
|:---:|:---:|
| PUT | Create |
| GET | Read |
| PUT | Update |
| DELETE | Delete  |
