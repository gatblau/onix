Feature: Delete Privilege
  As an API user
  I want to delete an existing privilege
  So that the user associated with a role cannot perform
    read, write and delete operations on resources in a partition.

  Scenario: Delete Privilege
    Given the privilege exists in the database
    When a DELETE HTTP request to the privilege endpoint is done
    Then the response code is 200
    Then the response has body
    Then the service responds with action "D"