# ox_partition Data Source  <img src="../../../docs/pics/ox.png" width="200" height="200" align="right">

Allow partition data to be fetched for use elsewhere in Terraform configuration.

More information about partitions can be found in the [Partition Resource](../resources/ox_partition.md) section.

## Example Usage

```hcl
data "ox_partition" "Logistics_Department_Partition" {
  key = "LOGISTICS_DEPT_PARTITION"
}
```

## Argument Reference

The data source requires the following arguments:

| Name | Use | Type |  Description |
|---|---|---|---|
| `key` | required | string | *The natural key that uniquely identifies the partition.* |

## Attribute Reference

The data source exports the following attributes:

| Name | Type |  Description |
|---|---|---|
| `name`| *required* | string | *The display name for the logical partition.* |
| `description`| *required* | string | *A meaningful description for the logical partition.* |
| `version` | optional | integer | *The version number of the partition for [optimistic concurrency control](https://en.wikipedia.org/wiki/Optimistic_concurrency_control) purposes. If specified, the entity can be written provided that the specified version number matches the one in the database. If no specified, optimistic locking is disabled.* |
| `created` | date & time | *The date and time the partition definition was first created.* |
| `updated` | date & time | *The date and time the partition was last updated.* |
| `changed_by` | string | *The user and role that last modified the partition.* |