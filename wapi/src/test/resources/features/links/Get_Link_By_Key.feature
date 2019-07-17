Feature: Get Item Link By Key
  As an API user
  I want to retrieve an item link
  So that a I can use the information in the CMDB

  Scenario: Get an Item Link using its natural key
    Given an item link natural key is known
    Given the link URL search by key is known
    Given the link exists in the database
    When a GET HTTP request to the Link with Key URL is done
    Then the response code is 200
    Then the response has body
    Then the reponse contains the requested link