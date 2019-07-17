Feature: Check that the service is alive
  As an API user
  I want to find if the service is running
  So that a liveliness probe can be set up for application management

  Scenario: Check the service is alive
    Given the live URL of the service is known
    When a get request to the live url is done
    Then the response code is 200