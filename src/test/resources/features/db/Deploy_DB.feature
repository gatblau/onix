Feature: Deploy Database
  As a system administrator
  I want to deploy a new Onix database when the readiness probe is called and the database does not exist
  So that the database is deployed automatically if not there

  Scenario: Deploy new database
    Given the database does not exist
    When the readiness probe is checked
    Then the database is deployed