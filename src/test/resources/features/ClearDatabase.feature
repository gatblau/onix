Feature: Clear the Configuration Management Database
  As an API user
  I want to clear all entries in the CMDB
  So that I can prepare tests to run on a clear database

  Scenario: Clear database
    Given the clear URL of the service is known
    When a delete request to the service is done
    Then there is not any error in the response
    Then the database is cleared