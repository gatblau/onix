Feature: Create Item Type Attribute
  As an API user
  I want to create a new Configuration Item type attribute
  So that it can be used to define an attribute of the item of this type

  Scenario: Create a new Configuration Item Type attribute
    Given a model exists in the database
    Given an item type exists in the database
    Given the item type attribute URL of the service with key is known
    Given the item type attribute natural key is known
    When an item type attribute PUT HTTP request with a JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the result contains no errors