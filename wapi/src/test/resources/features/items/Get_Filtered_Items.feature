Feature: Get Filtered Configuration Items
  As an API user
  I want to retrieve a list of configuration items by item type, item tag and date range
  So that a I can use the information in the CMDB

  Scenario: Get a list of Items that match item type, item tag and date range
    Given the item URL search with query parameters is known
    Given more than one item exist in the database
    Given the filtering config item type is known
    Given the filtering config item tag is known
    Given the filtering config item date range is known
    When a GET HTTP request to the Item uri is done with query parameters
    Then the response code is 200
    Then the response has body
    Then the response contains more than 2 items