Feature: Update Item Tree Tag
  As an API User
  I want to update name, description and/or label text of an existing tag
  So that I can improve the information the tag provides after it has been created.

  Scenario: Update tag data
    Given the URL of the tag update endpoint is known
    Given the tag already exists
    Given the item root key of the tag is known
    Given the current label of the tag is known
    Given a payload exists with the data required to update the tag
    When a tag update is requested
    Then the response code is 200
    Then the result contains no errors