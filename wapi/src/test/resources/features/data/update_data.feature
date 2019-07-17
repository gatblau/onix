Feature: Re-import data
  As an API user
  I want to update all items and links that define a particular configuration
  So that I facilitate the bulk update of large amount of configuration data at once.

  Scenario: Import data updating already imported data
    Given the URL of the item tree PUT endpoint is known
    Given the item tree exists in the database
    Given a json payload with update tree data exists
    When the update of the tree is requested
    Then the response code is 200
    Then the result list contains no errors
    Then the result list contained updated results