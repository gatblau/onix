Feature: Find a Role using its key
  As an Web API user
  I want to find a specific Role by key
  So that I can see its information and alternatively update or delete it later.

  Scenario: Find Role by Key
    Given a role natural key is known
    Given the role exists in the database
    When a request to GET the role is made
    Then the response has body
    Then the response contains the requested partition
    Then the response code is 200