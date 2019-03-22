Feature: Get meta models
  As an API user
  I want to get meta models in the system
  So that I can use them for presentation and reporting purposes.

  Scenario: Get all meta models
    Given the meta model URL of the service without key is known
    Given there are a few meta models in the system
    When a meta model GET HTTP request is done
    Then the response code is 200
    Then the response has body
    Then the response contains more than 1 meta models

  Scenario: Get meta model by key
    Given the meta model natural key is known
    Given the meta model URL of the service with key is known
    Given there are a few meta models in the system
    When a meta model GET HTTP request with key is done
    Then the response code is 200
    Then the response has body
    Then there is not any error in the response