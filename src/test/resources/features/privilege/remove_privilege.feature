Feature: Remove Privilege
  As an API user
  I want to remove an existing privilege from a role
  So that the user associated with the role cannot perform
    read, write and delete operations on resources.

  Scenario: Remove Privilege
    Given the privilege exists in the database
    When a DELETE HTTP request to the privilege endpoint is done
    Then the response code is 200
    Then the response has body
    Then the service responds with action "D"