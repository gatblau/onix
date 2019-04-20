Feature: Update Configuration Item
  As an API user
  I want to update an existing configuration item
  So that a I can record changes to the item as a result of deployment automation

  Scenario: Update an existing Item using a JSON payload
    Given a configuration item natural key is known
    Given a model exists in the database
    Given an item type exists in the database
    Given the item exist in the database
    Given the item URL search by key is known
    When a PUT HTTP request with an updated JSON payload is done
    Then the response code is 200
    Then the response has body
    Then the service responds with action "U"