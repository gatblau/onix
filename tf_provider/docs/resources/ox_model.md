# ox_model Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Creates, updates or destroys a data model definition.

A data model defines the list of item types, link types and link rules that are used for a particular configuration storage requirement.

A data model has to be defined before configuration items acan be created and linked.

## Example Usage

```hcl
resource "ox_model" "Test_Model" {
  key         = "test_model"
  name        = "Test Model"
  description = "Test Model Description"
  managed     = false
}
```

## Argument Reference

The following arguments can be passed to a configuration item:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the model.* |
| `name`| required | string | *The display name for the model.* |
| `description`| required | string | *A meaningful description for the model.* |
| `partition`| optional | string | *The logical access partition the model is in. If not specified, the default reference partition (REF) is used.* |
| `managed` | optional | boolean | *A flag that informs whether the model is managed by an external application process. The default value is FALSE.* |