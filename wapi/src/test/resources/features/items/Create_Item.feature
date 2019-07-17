Feature: Create Configuration Item
  As an API user
  I want to create a configuration item
  So that a I can record required information in the CMDB

  Scenario: Create an Item using a JSON payload
    Given a configuration item natural key is known
    Given the item does not exist in the database
    Given a model exists in the database
    Given an item type exists in the database
    Given a configuration item natural key is known
    Given the item URL search by key is known
    When a PUT HTTP request with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"