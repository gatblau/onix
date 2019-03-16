Feature: Delete Item Tree Snapshot
  As an API User
  I want to delete all snapshots for a given root item
  So that I can keep the database clean if the tag is not required.

  Scenario: Delete tag
    Given the URL of the tag delete all endpoint is known
    Given there are more than one tags in the database
    Given the item root key of the tag is known
    When a tag delete all is requested
    Then the response code is 200
    Then the result contains no errors