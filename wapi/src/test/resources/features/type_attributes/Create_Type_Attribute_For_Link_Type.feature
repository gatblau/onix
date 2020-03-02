Feature: Create Type Attribute for an Link Type
  As an API user
  I want to create a new type attribute for a link type
  So that client applications understand the allowed properties for a link type

  Scenario: Create a new Type Attribute for a Link Type
    Given the key of the link type is known
    Given the key of the type attribute for the link type is known
    Given the link type attribute does not exist in the database
    Given a model exists in the database
    Given an link type exists in the database
    Given the link type attribute natural key is known
    When a PUT HTTP request with a JSON payload is done for an attribute of a link type
    Then the response code is 201
    Then the response has body
    Then the result contains no errors