Feature: Delete Partition
  As an API user
  I want to delete an existing logical partition
  So that the partitions not used do not use space in the CMDB.

  Scenario: Delete Partition
    Given a partition natural key is known
    Given the partition exists in the database
    Given the partition DELETE URL by key is known
    When a DELETE HTTP request to the partition resource with an item key is done
    Then the response has body
    Then the service responds with action "D"