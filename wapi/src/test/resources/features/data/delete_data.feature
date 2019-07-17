Feature: Delete an existing Item and children
  As an API user
  I want to delete all items and links that define a particular configuration
  So that I facilitate the removal of large amount of configuration data at once.

  Scenario: Delete item tree
    Given the URL of the item tree DELETE endpoint is known
    Given the item tree exists in the database
    Given the item key of the tree root item is known
    Given there are not any tag for the root item
    When the deletion of the tree is requested
    Then the response code is 200
    Then the result contains no errors