Feature: Delete Link Between two Configuration Items
  As an API user
  I want to delete a link between two existing configuration items
  So that a I can remove the association in the CMDB

  Scenario: Delete a link
    Given the natural key for the link is known
    Given the link URL of the service is known
    When a DELETE Link request is done
    Then there is not any error in the response