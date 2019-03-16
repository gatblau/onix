Feature: Get Item Tree Tags
  As an API User
  I want to get all existing tags for a given root item
  So that I can decide which tag to select based on their given labels.

  Scenario: Get tags
    Given the URL of the tag get endpoint is known
    Given there are tags for a given item in the database
    Given the item root key of the tags is known
    When a tag list for an item is requested
    Then the response code is 200
    Then the response contains more than 0 tags