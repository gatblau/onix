Feature: Create privilege
  As an Admin user
  I want to create a privilege
  So that the user associated with a specific role can perform
    CRUD operations on resources in a specific partition.

  Scenario: Create privilege
    Given the partition exists in the database
    Given the role exists in the database
    Given the privilege does not exist in the database
    When a PUT HTTP request to the privilege endpoint with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"