# ox_role Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_role resource, creates, updates or destroys access control roles.

A role is used to attach [privileges](ox_privilege.md) to [logical partitions](ox_partition.md).

When a user logs in the Onix Web API, a role is assigned to the user session after the user is authenticated. The role is passed to the database so that any CRUD operation is approved based on the [privileges](ox_privilege.md) attached to the role.

## Example usage

```hcl
resource "ox_role" "Logistics_Admin_Role" {
  key         = "LOGISTICS_ADMIN_ROLE"
  name        = "Logistics Department Administrator"
  description = "Role held by an Administrator of the Logistics Department."
  level       = 1
}
```

## Argument reference

The following arguments can be passed to a role:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | *required* | string | *The natural key that uniquely identifies the role.* |
| `name`| *required* | string | *The display name for the role.* |
| `description`| *required* | string | *A meaningful description for the role.* |
| `level` | optional | boolean | *The administration level determines the privilege the role has to create, update, read and delete partitions and other roles. The default value is 0. </br> Level 0 roles cannot modify or read neither partitions nor roles data. </br> Level 1 roles can only modify / read partitions or roles data for the ones which are owned by the role.</br> Level 2 roles can modify / read any partition or role data.* |
| `version` | optional | integer | *The version number of the partition for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |

## Key dependencies

A role has no other dependencies.

![Role](../pics/role.png)

## Related resources

- Role **has** [Privilege](ox_privilege.md)
