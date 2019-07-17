Feature: Import data
  As an API user
  I want to create all items and links that define a particular configuration
  So that I facilitate the import of large amount of configuration data at once.

  Scenario: Import data graph
    Given the URL of the item tree PUT endpoint is known
    Given the item tree does not exist in the database
    Given a json payload with tree data exists
    When the creation of the tree is requested
    Then the response code is 200
    Then the result list contains no errors