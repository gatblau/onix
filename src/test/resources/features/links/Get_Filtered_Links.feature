Feature: Get Filtered Item Links
  As an API user
  I want to retrieve a list of item links by various filters
  So that a I can use the information in the CMDB

  Scenario: Get a list of Links that match item type, item tag and date range
    Given the link URL search with query parameters is known
    Given more than one link exist in the database
    Given the filtering link type is known
    Given the filtering link tag is known
    Given the filtering link date range is known
    When a GET HTTP request to the Link uri is done with query parameters
    Then the response code is 200
    Then the response has body
    Then the response contains more than 2 links