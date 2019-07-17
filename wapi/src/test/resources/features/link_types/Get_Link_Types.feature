Feature: Find Link Types
  As an Web API user
  I want to find Link Types
  So that I can see their configuration for administration purposes.

  Scenario: Find All Link Types
    Given there are pre-existing Link types in the database
    Given the link type URL of the service with no query parameters exist
    When a request to GET a list of link types is done
    Then the response has body
    Then the response contains more than 2 link types
    Then the response code is 200