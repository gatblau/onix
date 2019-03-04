Feature: Get Ansible Inventory

  Scenario: Get Inventory using Key
    Given the inventory key is known
    Given the inventory snapshot label is known
    Given the snapshot already exists
    Given the inventory exists in the database
    Given the URL of the inventory finder endpoint is known
    When an HTTP GET to the inventory GET endpoint is made using its key
    Then the response code is 200
