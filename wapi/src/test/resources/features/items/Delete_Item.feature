Feature: Delete Configuration Item
  As an API user
  I want to delete an existing configuration item
  So that a I can record the fact the item has been decommissioned.

  Scenario: Delete Configuration Item
    Given the meta model natural key is known
    Given the meta model exists in the database
    Given the item type exists in the database
    Given the item exists in the database
    Given the item URL search by key is known
    When a DELETE HTTP request with an item key is done
    Then the response has body
    Then the service responds with action "D"