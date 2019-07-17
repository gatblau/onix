Feature: Delete All Link Rules
  As an API user
  I want to delete all Link Rules
  So that I can start configuring the CMDB from an empty list of rules

  Scenario: Delete all Link Rules
    Given the link rule URL of the service is known
    When an link rule DELETE HTTP request is done
    Then there is not any error in the response