Feature: Find a Partition using its key
  As an Web API user
  I want to find a specific Partition by key
  So that I can see its information and alternatively update or delete it later.

  Scenario: Find Partition by Key
    Given the partition natural key is known
    Given the partition is in the database
    When a request to GET the partition is made
    Then the response has body
    Then the response contains the requested partition
    Then the response code is 200