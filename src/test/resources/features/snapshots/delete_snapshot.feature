Feature: Delete Item Tree Snapshot
  As an API User
  I want to delete an existing snapshot
  So that I can keep the database clean if the snapshot is not required.

  Scenario: Delete snapshot
    Given the URL of the snapshot delete endpoint is known
    Given the snapshot already exists
    Given the item root key of the snapshot is known
    Given the current label of the snapshot is known
    When a snapshot delete is requested
    Then the response code is 200
    Then the result contains no errors