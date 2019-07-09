Feature: Delete a Link Rule
  As an API user
  I want to delete a specific Link Rule
  So that I can remove link constraints that are no longer needed

  Scenario: Delete a Link Rule
    Given the meta model exists in the database
    Given the item type exists in the database
    Given the link type exists in the database
    Given the link rule URL of the service is known
    Given the link rule natural key is known
    Given the link rule exists in the database
    When a DELETE HTTP request with a link rule key is done
    Then there is not any error in the response