Feature: Delete Item Tree Tag
  As an API User
  I want to delete an existing tag
  So that I can keep the database clean if the tag is not required.

  Scenario: Delete tag
    Given the URL of the tag delete endpoint is known
    Given the tag already exists
    Given the item root key of the tag is known
    Given the current label of the tag is known
    When a tag delete is requested
    Then the response code is 200
    Then the result contains no errors