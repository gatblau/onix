Feature: Update Item Tree Snapshot
  As an API User
  I want to update name, description and/or label text of an existing snapshot
  So that I can improve the information the snapshot provides after it has been created.

  Scenario: Update snapshot data
    Given the URL of the snapshot update endpoint is known
    Given the snapshot already exists
    Given the item root key of the snapshot is known
    Given the current label of the snapshot is known
    Given a payload exists with the data required to update the snapshot
    When a snapshot update is requested
    Then the response code is 200
    Then the result contains no errors