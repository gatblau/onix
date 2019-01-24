Feature: Delete All Item types
  As an API user
  I want to delete all Configuration Item types
  So that I can start configuring the CMDB from an empty list of types

  Scenario: Delete all Configuration Item types
    Given the item type URL of the service is known
    When an item type DELETE HTTP request is done
    Then there is not any error in the response