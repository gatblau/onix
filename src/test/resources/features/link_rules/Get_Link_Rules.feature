Feature: Find Link Rules
  As an Web API user
  I want to find Link Rules
  So that I can see their configuration for administration purposes.

  Scenario: Find All Link Rules
    Given: there are pre-existing link rules in the database
    Given the link rule URL of the service with no query parameters exist
    When a request to GET a list of link rules is done
    Then the response has body
    Then the response contains 2 link rules
    Then the response code is 200