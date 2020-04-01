Feature: Find Link Types
  As an Web API user
  I want to find Link Types
  So that I can see their configuration for administration purposes.

  Scenario: Find All Link Types
    Given the meta model natural key is known
    Given the meta model exists in the database
    Given there are item types in the database
    Given there are link types in the database
    Given there are link type attributes for the link types in the database
    Given the link type attribute URL exist
    When a request to GET a list of link type attributes is done
    Then the response has body
    Then the response contains more than 1 link type attributes
    Then the response code is 200