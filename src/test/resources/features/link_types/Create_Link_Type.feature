Feature: Create a new Configuration Item type
  As an API user
  I want to create a new Link Item type
  So that links can be of the created link type

  Scenario: Create a new Link Type
    Given the link type does not exist in the database
    Given the link type URL of the service with key is known
    Given the link type natural key is known
    Given a json payload with new link type information exists
    When a link type PUT HTTP request with a JSON payload is done
    Then the response code is 200
    Then the response has body