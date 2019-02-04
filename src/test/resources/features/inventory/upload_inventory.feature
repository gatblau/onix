Feature: Upload Ansible Inventory
  As an API user
  I want to upload an Ansible inventory
  So that I can store the inventory to be used by Ansible later.

  Scenario: Upload new inventory
    Given an inventory file exists
    Given the inventory key is known
    Given the inventory upload URL is known
    When an HTTP PUT request with the inventory payload is executed
    Then the inventory config item is created
    Then the host group config items are created
    Then the host config items are created