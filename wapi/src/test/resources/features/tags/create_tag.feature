Feature: Create Item Tree Tag
  As an API User
  I want to take a tag of all items and links that are connected to a root item
  So that I can retrieve those items and links as they were at a particular point in time.

  Scenario: Take a tag
    Given the URL of the tag create endpoint is known
    Given there are items linked to the root item in the database
    Given a payload exists with the data required to create the tag
    Given the tag does not already exist
    When a tag creation is requested
    Then the response code is 201
    Then the result contains no errors