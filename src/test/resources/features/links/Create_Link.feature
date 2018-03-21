Feature: Create Link Between two Configuration Items
  As an API user
  I want to create a link between two existing configuration items
  So that a I can record and association in the CMDB

  Scenario: Create a link using a JSON payload
    Given the natural keys for two configuration items are known
    Given the link URL of the service is known
    Given a link to the two configuration items does not exist in the database
    Given a json payload with new link information exists
    When a PUT HTTP request with a JSON payload is done
    Then the response code is 200
    Then the response has body
    Then the service responds with action "CREATED"