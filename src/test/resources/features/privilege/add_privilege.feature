Feature: Add privilege to Role
  As an Admin user
  I want to add privileges to a role
  So that the user associated with the role can perform
    read, write and delete operations on resources.

  Scenario: Add privilege to role
    Given the partition exists in the database
    Given the role exists in the database
    Given the privilege does not exist in the database
    When a PUT HTTP request to the privilege endpoint with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"