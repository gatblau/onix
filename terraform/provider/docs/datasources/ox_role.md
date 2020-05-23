# ox_partition Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow role data to be fetched for use elsewhere in Terraform configuration.

More information about partitions can be found in the [Role Resource](../resources/ox_role.md) section.

## Example Usage

```hcl
data "ox_role" "Logistics_Admin_Role" {
  key = "LOGISTICS_ADMIN_ROLE"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the role.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| string | *The display name for the role.* |
| `description`| string | *A meaningful description for the role.* |
| `level` | boolean | *The administration level determines the privilege the role has to create, update, read and delete partitions and other roles. The default value is 0. </br> Level 0 roles cannot modify or read neither partitions nor roles data. </br> Level 1 roles can only modify / read partitions or roles data for the ones which are owned by the role.</br> Level 2 roles can modify / read any partition or role data.* |
| `version` | integer | *The version number of the role for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the role definition was first created.* |
| `updated` | date & time | *The date and time the role was last updated.* |
| `changed_by` | string | *The user and role that last modified the role.* |