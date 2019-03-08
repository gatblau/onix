Feature: Get Item Tree Snapshots
  As an API User
  I want to get all existing snapshots for a given root item
  So that I can decide which snapshot to select based on their given labels.

  Scenario: Get snapshots
    Given the URL of the snapshot get endpoint is known
    Given there are snapshots for a given item in the database
    Given the item root key of the snapshots is known
    When a snapshot list for an item is requested
    Then the response code is 200
    Then the response contains more than 0 snapshots