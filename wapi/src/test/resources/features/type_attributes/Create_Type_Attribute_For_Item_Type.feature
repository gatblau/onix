Feature: Create Type Attribute for an Item Type
  As an API user
  I want to create a new type attribute for an item type
  So that client applications understand the allowed properties for an item type

  Scenario: Create a new Type Attribute for an Item Type
    Given the key of the item type is known
    Given the key of the type attribute for the item type is known
    Given the item type attribute does not exist in the database
    Given a model exists in the database
    Given an item type exists in the database
    Given the item type attribute natural key is known
    When a PUT HTTP request with a JSON payload is done for an attribute of an item type
    Then the response code is 201
    Then the response has body
    Then the result contains no errors