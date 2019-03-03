Feature: Create Item Tree Snapshot
  As an API User
  I want to take a snapshot of all items and links that are connected to a root item
  So that I can retrieve those items and links as they were at a particular point in time.

  Scenario: Take a snapshot
    Given the URL of the snapshot create endpoint is known
    Given there are items linked to the root item in the database
    Given a payload exists with the data required to create the snapshot
    Given the snapshot does not already exist
    When a snapshot creation is requested
    Then the response code is 200
    Then the result contains no errors