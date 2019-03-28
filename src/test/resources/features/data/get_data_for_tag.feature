Feature: Get Item and children for a specified Tag
  As an API user
  I want to retrieve all items and links that are associated with a root item
    that has a tag identified by a label
  So that I can visualise the tree information for management purposes.

  Scenario: Get tag item tree
    Given the URL of the item tree get endpoint is known
    Given the tag already exists
    Given the item root key of the tag is known
    Given the current label of the tag is known
    When a tag tree retrieval for the tag is requested
    Then the response code is 200
    Then the result contains the tree items and links