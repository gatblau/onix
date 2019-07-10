provider "ox" {
  uri  = "http://localhost:8080"
  user = "admin"
  pwd  = "0n1x"
}

resource "ox_model" "ansible" {
  key         = "ANSIBLE_INVENTORY"
  name        = "Ansible Inventory"
  description = "Provides the required structure to store Ansible inventory information."
}

resource "ox_item_type" "ansible_inventory" {
  key         = "ANSIBLE_INVENTORY"
  name        = "Ansible Inventory"
  description = "The Ansible inventory."
  model_key   = "${ox_model.ansible.key}"
}

resource "ox_item_type" "ansible_host_super_set" {
  key         = "ANSIBLE_HOST_SUPER_SET"
  name        = "Ansible Host Super Set"
  description = "A group of Ansible groups."
  model_key   = "${ox_model.ansible.key}"
}

resource "ox_item_type" "ansible_host_group" {
  key         = "ANSIBLE_HOST_GROUP"
  name        = "Ansible Host Group"
  description = "A group of hosts."
  model_key   = "${ox_model.ansible.key}"
}

resource "ox_item_type" "ansible_host" {
  key         = "ansible_host"
  name        = "Ansible Host"
  description = ""
  model_key   = "${ox_model.ansible.key}"
}

resource "ox_link_type" "ansible_inventory_link" {
  key         = "ANSIBLE_INVENTORY_LINK"
  name        = "Ansible Inventory Link"
  description = "Links Ansible inventory items."
  model_key   = "${ox_model.ansible.key}"
}

resource "ox_link_rule" "ansible_inventory_to_ansible_host_super_set" {
  key                 = "ANSIBLE_INVENTORY->ANSIBLE_HOST_SUPER_SET"
  name                = "Ansible Inventory to Host Super Set rule"
  description         = "Allow linking Inventory with Host Super Set items"
  link_type_key       = "${ox_link_type.ansible_inventory_link.key}"
  start_item_type_key = "${ox_item_type.ansible_inventory.key}"
  end_item_type_key   = "${ox_item_type.ansible_host_super_set.key}"
}

resource "ox_link_rule" "ansible_inventory_to_ansible_host_group" {
  key                 = "ANSIBLE_INVENTORY->ANSIBLE_HOST_GROUP"
  name                = "Ansible Inventory to Host Group rule"
  description         = "Allow linking Inventory with Host Group items"
  link_type_key       = "${ox_link_type.ansible_inventory_link.key}"
  start_item_type_key = "${ox_item_type.ansible_inventory.key}"
  end_item_type_key   = "${ox_item_type.ansible_host_group.key}"
}

resource "ox_link_rule" "ansible_host_super_set_to_ansible_host_group" {
  key                 = "ANSIBLE_HOST_SUPER_SET->ANSIBLE_HOST_GROUP"
  name                = "Ansible Host Super Set to Host Group rule"
  description         = "Allow linking Ansible Host Super Sets with Host Groups set items"
  link_type_key       = "${ox_link_type.ansible_inventory_link.key}"
  start_item_type_key = "${ox_item_type.ansible_host_super_set.key}"
  end_item_type_key   = "${ox_item_type.ansible_host_group.key}"
}

resource "ox_link_rule" "ansible_host_group_to_ansible_host" {
  key                 = "ANSIBLE_HOST_GROUP->ANSIBLE_HOST"
  name                = "Ansible Host Group to Host rule"
  description         = "Allow linking Ansible Host Group with Host items"
  link_type_key       = "${ox_link_type.ansible_inventory_link.key}"
  start_item_type_key = "${ox_item_type.ansible_host_group.key}"
  end_item_type_key   = "${ox_item_type.ansible_host.key}"
}
