Feature: Create Meta Model
  As an API user
  I want to create a new meta model
  So that I can group item types and link types into models
   for retrieval and data management purposes

  Scenario: Create a new Meta Model
    Given the meta model natural key is known
    Given the meta model does not exist in the database
    Given the meta model URL of the service with key is known
    When a meta model PUT HTTP request with a JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the result contains no errors