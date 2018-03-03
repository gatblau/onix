Feature: Create a new Configuration Item type
  As an API user
  I want to create a new Configuration Item type
  So that I can assign it to Configuration Items

  Scenario: Create a new Configuration Item type
    Given the item type does not exist in the database
    Given the item type URL of the service is known
    Given a json payload with new item type information exists
    When a POST HTTP request with a JSON payload is done
    Then the response code is 200
    Then the response has body