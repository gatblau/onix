Feature: Find Link Type by Key
  As an Web API user
  I want to find a specific Link Type
  So that I can see its configuration for administration purposes.

  Scenario: Find Link Type by Key
    Given the link type exists in the database
    Given the link type URL of the service with key is known
    Given the link type natural key is known
    When a link type GET HTTP request with the key is done
    Then the response has body
    Then the response code is 200
    Then the response contains the link type