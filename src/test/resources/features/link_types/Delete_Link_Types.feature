Feature: Delete All Item types
  As an API user
  I want to delete all Link types
  So that I can start configuring the CMDB from an empty list of types

  Scenario: Delete all Link types
    Given the link type URL of the service is known
    When a link type DELETE HTTP request is done
    Then there is not any error in the response