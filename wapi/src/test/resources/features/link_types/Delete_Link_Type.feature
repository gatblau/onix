Feature: Delete an Item type
  As an API user
  I want to delete a specific Link type
  So that I can avoid storing types that are no longer needed

  Scenario: Delete a Link type
    Given the link type URL of the service is known
    Given the link type natural key is known
    Given the link type exists in the database
    When a DELETE HTTP request with a link type key is done
    Then there is not any error in the response