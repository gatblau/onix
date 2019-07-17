Feature: Delete an Item type
  As an API user
  I want to delete a specific Configuration Item type
  So that I can avoid storing types that are no longer needed

  Scenario: Delete a Configuration Item type
    Given the item type URL of the service is known
    Given the item type natural key is known
    Given the item type exists in the database
    When a DELETE HTTP request with an item type key is done
    Then there is not any error in the response