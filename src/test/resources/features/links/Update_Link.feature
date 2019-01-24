Feature: Update Link
  As an API user
  I want to update an existing link
  So that a I can update the properties of the link as they change

  Scenario: Update an existing Link using a JSON payload
    Given the natural key for the link is known
    Given the configuration items to be linked exist in the database
    Given the link between the two items exists in the database
    Given the link URL of the service is known
    When a PUT HTTP request with an updated JSON payload is done
    Then the response code is 200
    Then the response has body
    Then the service responds with action "U"