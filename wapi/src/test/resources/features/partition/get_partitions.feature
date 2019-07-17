Feature: Find Partitions
  As an Web API user
  I want to find Partitions
  So that I can see their configuration for administration purposes.

  Scenario: Find All Partitions
    Given there are multiple partitions in the database
    Given the partition URL of the service with no query parameters exist
    When a request to GET a list of partitions is done
    Then the response has body
    Then the response contains more than 1 partitions
    Then the response code is 200