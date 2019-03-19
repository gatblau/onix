Feature: Get Configuration Item Metadata
  As an API user
  I want to retrieve a configuration item metadata
  So that a I can consume it directly from client applications.

  Scenario: Get an Item Unfiltered Metadata
    Given a configuration item natural key is known
    Given the item metadata URL get by key is known
    Given the item exists in the database
    When a GET HTTP request to the Item Metadata endpoint is done
    Then the response code is 200
    Then the response has body
    Then the response contains the requested metadata