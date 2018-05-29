Feature: Get all links for a configuration item
  As an API user
  I want to retrieve all links for a specified configuration item
  So that a I can understand any relation between the item and other items.

  Scenario: Get all item links
    Given a configuration item natural key is known
    Given the links by item URL of the service is known
    Given two items exist in the database
    Given two links between the two configuration items exist in the database
    When a GET HTTP request to the Link by Item resource is done
    Then the response code is 200
    Then the response has body
    Then the response contains 2 links