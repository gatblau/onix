Feature: Create Link Between two Configuration Items
  As an API user
  I want to create a link between two existing configuration items
  So that a I can record and association in the CMDB

  Scenario: Create a link using a JSON payload
    Given the meta model natural key is known
    Given the meta model exists in the database
    Given the item type natural key is known
    Given the item type exists in the database
    Given the natural key for the link is known
    Given the link URL of the service is known
    Given the link type exists in the database
    Given the link rule exists in the database
    Given the configuration items to be linked exist in the database
    Given the link type natural key is known
    Given the link type exists in the database
    Given the link rule exists in the database
    Given a link to the two configuration items does not exist in the database
    Given a json payload with new link information exists
    When a PUT HTTP request to the Link resource is done with a JSON payload
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"
    Then the result contains no errors