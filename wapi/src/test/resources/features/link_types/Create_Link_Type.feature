Feature: Create a new Link Type
  As an API user
  I want to create a new Link Type
  So that link instances can be associated to the link type

  Scenario: Create a new Link Type
    Given the link type does not exist in the database
    Given a model exists in the database
    Given the link type URL of the service with key is known
    Given the link type natural key is known
    When a link type PUT HTTP request with a JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the result contains no errors