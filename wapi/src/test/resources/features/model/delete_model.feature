Feature: Delete Model
  As an API user
  I want to delete an existing meta model
  So that I can remove any item types and link types
    associated with a no longer used meta model

  Scenario: Delete existing Meta Model
    Given the meta model natural key is known
    Given the meta model exists in the database
    Given there are not any items associated with the model
    Given there are not any item types associated with the model
    Given there are not any link types associated with the model
    Given the meta model URL of the service with key is known
    When a meta model DELETE HTTP request with key is done
    Then the response code is 200
    Then the response has body
    Then the result contains no errors
    Then the service responds with action "D"