Feature: Find Roles
  As an Web API user
  I want to find a list of all roles
  So that I can see their configuration for administration purposes.

  Scenario: Find All Roles
    Given there are multiple roles in the database
    Given the role URL of the service with no query parameters exist
    When a request to GET a list of roles is done
    Then the response has body
    Then the response contains more than 1 roles
    Then the response code is 200