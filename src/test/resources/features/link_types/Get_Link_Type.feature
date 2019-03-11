Feature: Find Link Type by Key
  As an Web API user
  I want to find a specific Link Type
  So that I can see its configuration for administration purposes.

  Scenario: Find Link Type by Key
    Given there is a link type in the database
    Given the link type URL GET exists
    When a request to GET the link type is done
    Then the response has body
    Then the response code is 200
    Then the response contains the link type