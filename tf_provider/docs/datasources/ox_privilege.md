# ox_privilege Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow privilege data to be fetched for use elsewhere in Terraform configuration.

More information about partitions can be found in the [Privilege Resource](../resources/ox_privilege.md) section.

## Example Usage

```hcl
data "ox_privilege" "Logistics_Department_Reader_Privilege" {
  key = "PRIVILEGE_LOGISTICS_DEPT_READER"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the privilege.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name` | string | *The display name for the logical partition.* |
| `description` | string | *A meaningful description for the logical partition.* |
| `version` | integer | *The version number of the privilege for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the privilege definition was first created.* |
| `updated` | date & time | *The date and time the privilege was last updated.* |
| `changed_by` | string | *The user and role that last modified the privilege.* |