Feature: Get Configuration Item Metadata
  As an API user
  I want to retrieve a configuration item metadata
  So that a I can consume it directly from client applications.

  Scenario: Get an Item Unfiltered Metadata
    Given a model exists in the database
    Given a configuration item natural key is known
    Given a metadata filter key is known
    Given the item metadata URL GET with filter is known
    Given an item type with filter data exists in the database
    Given the item with metadata exists in the database
    When a GET HTTP request to the Item Metadata endpoint with filter is done
    Then the response code is 200
    Then the response has body
    Then the response contains the requested metadata
    Then there is not any error in the response