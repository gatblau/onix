Feature: Get Configuration Item By Key
  As an API user
  I want to retrieve a configuration item
  So that a I can use the information in the CMDB

  Scenario: Get an Item using its natural key
    Given a configuration item natural key is known
    Given the item URL search by key is known
    Given the item exists in the database
    When a GET HTTP request to the Item uri is done
    Then the response code is 200
    Then the response has body
    Then the reponse contains the requested item