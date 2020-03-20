Feature: Create a new Link Rule
  As an API user
  I want to create a new Link Rule
  So that I can restrict the types of items that can be linked

  Scenario: Create a new Link Rule
    Given the link rule does not exist in the database
    Given the meta model exists in the database
    Given the item types to and from exists in the database
    Given the link type exists in the database
    Given the link rule URL of the service with key is known
    Given the link rule natural key is known
    When a link rule PUT HTTP request with a JSON payload is done
    Then the response code is 201
    Then the response has body