Feature: Create User
  As an API user
  I want to create a new user account
  So that the account can be used for accessing the Web API.

  Scenario: Create a User account using a JSON payload
    Given a user account natural key is known
    Given the user account does not exist in the database
    Given the user account PUT URL by key is known
    When a PUT HTTP request to the user account endpoint with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"