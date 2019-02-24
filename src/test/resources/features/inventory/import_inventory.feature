Feature: Import Ansible Inventory
  As an API user
  I want to upload an Ansible inventory
  So that I can store the inventory to be used by Ansible later.

  Scenario: Upload new inventory
    Given an inventory file exists
    Given the inventory key is known
    Given the inventory upload URL is known
    When an HTTP PUT request with the inventory payload is executed
    Then there is not any error in the response