Feature: Create Node
  As an API user
  I want to create a configuration item
  So that a I can record required information in the CMDB

  Scenario: Create a Node with JSON payload
    Given a node key is known
    Given the node does not exist in the database
    Given the create URL of the service is known
    Given a json payload with node information exists
    When a create JSON request to the service is done
    Then the response code is 200
    Then the response has body
    Then the service responds with the create node Id