Feature: Create Role
  As an API user
  I want to create a new role
  So that a I can allocate privileges to the role for access control purposes.

  Scenario: Create a Role using a JSON payload
    Given a role natural key is known
    Given the role does not exist in the database
    Given the role PUT URL by key is known
    When a PUT HTTP request to the role endpoint with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"