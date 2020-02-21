Feature: Find Item Types
  As an Web API user
  I want to find Item Types
  So that I can see their configuration for administration purposes.

  Scenario: Find All Item Types
    Given the meta model natural key is known
    Given the meta model exists in the database
    Given there are item types in the database
    Given there are item type attributes for the item types in the database
    Given the item type attribute URL exist
    When a request to GET a list of item type attributes is done
    Then the response has body
    Then the response contains more than 1 item type attributes
    Then the response code is 200