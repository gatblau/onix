Feature: Check that the service is alive
  As an API user
  I want to find if the service is running
  So that a readiness probe can be set up for application management

  Scenario: Check the service is alive
    Given the base URL of the service is known
    When a get request to the service is done
    Then the service responds with description and version number