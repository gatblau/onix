Feature: Create a new Configuration Item type
  As an API user
  I want to create a new Configuration Item type
  So that I can assign it to Configuration Items

  Scenario: Create a new Configuration Item type
    Given a model exists in the database
    Given the item type does not exist in the database
    Given the item does not exist in the database
    Given the item type URL of the service with key is known
    Given the item type natural key is known
    When an item type PUT HTTP request with a JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the result contains no errors