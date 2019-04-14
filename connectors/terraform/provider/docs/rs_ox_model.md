# ox_model resource

The __ox_model resource__ allows the creation of new Onix modelKey definition, the update of existing modelKey definitions and the deletion of modelKey definitions.

## Example Usage

```hcl-terraform
resource "ox_model" "Test_Model" {
  key         = "test_model"
  name        = "Test Model"
  description = "Test Model Description"
}
```
## Argument Reference

The following passed-in arguments are supported:

| Argument | Description | Mandatory |
|---|---|---|
| __key__| the modelKey natural key | no |

__Note:__ If the __key__ argument is not specified, then all models are retrieved.

## Attribute Reference

The following returned attributes are supported:

| Attribute | Description |
|---|---|
| __name__| the modelKey human readable name. |
| __description__| the modelKey description. |