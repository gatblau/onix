Feature: Find Privileges
  As an Web API user
  I want to find all Privileges for a particular Role
  So that I can see the Role can do.

  Scenario: Find All Privileges by Role
    Given there are multiple privileges for a role in the database
    When a request to GET a list of privileges by role is done
    Then the response has body
    Then the response contains more than 1 privileges
    Then the response code is 200