Feature: Create Partition
  As an API user
  I want to create a new logical partition
  So that a I can segregate data in the CMDB for access control purposes

  Scenario: Create a Partition using a JSON payload
    Given a partition natural key is known
    Given the partition does not exist in the database
    Given the partition PUT URL by key is known
    When a PUT HTTP request to the partition endpoint with a new JSON payload is done
    Then the response code is 201
    Then the response has body
    Then the service responds with action "I"