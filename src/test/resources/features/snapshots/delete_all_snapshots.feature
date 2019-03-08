Feature: Delete Item Tree Snapshot
  As an API User
  I want to delete all snapshots for a given root item
  So that I can keep the database clean if the snapshot is not required.

  Scenario: Delete snapshot
    Given the URL of the snapshot delete all endpoint is known
    Given there are more than one snapshots in the database
    Given the item root key of the snapshot is known
    When a snapshot delete all is requested
    Then the response code is 200
    Then the result contains no errors