# ox_partition Resource <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

The ox_partition resource, creates, updates or destroys logical partitions.

A logical partition is a special tag which applies to either models (for reference data) and items (for instance data). Once models or items are tagged with a partition, a user role can be granted privilege to create, read or delete information within the partition.

For example, suppose that there are configuration data that is relevant to a part of an organization. A partition can be created to provide role based access control to any data created, updated, read or deleted within the partition.

A partition has no dependencies with any other Onix entity.

## Example usage

```hcl
resource "ox_partition" "Logistics_Department_Partition" {
  key         = "LOGISTICS_DEPT_PARTITION"
  name        = "Logistics Department Partition"
  description = "Partition for resources used by the logistic department."
  managed     = false
}
```

## Argument reference

The following arguments can be passed to a logical partition:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | *required* | string | *The natural key that uniquely identifies the logical partition.* |
| `name`| *required* | string | *The display name for the logical partition.* |
| `description`| *required* | string | *A meaningful description for the logical partition.* |
| `version` | optional | integer | *The version number of the partition for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |

## Key dependencies

A logical partition has no other dependencies.

![Partition](../pics/partition.png)

## Related resources

- [Model](ox_model.md) **is in** Partition
- [Item](ox_item.md) **is in** Partition
- [Privilege](ox_privilege.md) **allows to create, read or delete data in a** Partition
