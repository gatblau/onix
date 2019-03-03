Feature: Get Item Tree for a specified Snapshot
  As an API user
  I want to retrieve all items and links that are associated with a root item
    that has a snapshot identified by a label
  So that I can visualise the tree information for management purposes.

  Scenario: Get snapshot item tree
    Given the URL of the item tree get endpoint is known
    Given the snapshot already exists
    Given the item root key of the snapshot is known
    Given the current label of the snapshot is known
    When a snapshot tree retrieval for the snapshot is requested
    Then the response code is 200
    Then the result contains the tree items and links