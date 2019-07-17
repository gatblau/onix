Feature: Delete Role
  As an API Admin user
  I want to delete an existing role
  So that the database does not contain dated roles.

  Scenario: Delete Role
    Given a role natural key is known
    Given the role exists in the database
    Given the role DELETE URL by key is known
    When a DELETE HTTP request to the role resource with an item key is done
    Then the response has body
    Then the service responds with action "D"