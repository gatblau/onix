terraform {
  backend "local" {
    path = "terraform.d/terraform.tfstate"
  }
}

provider "onix" {
  uri = "http://localhost:8080"
  user = "admin"
  pwd = "0n1x"
}

resource "ox_model" "Test_Model" {
  key = "test_model"
  name = "Test Model"
  description = "Test Model Description"
}

resource "ox_item_type" "Test_Item_Type" {
  key = "test_item_type"
  name = "Test Item Type"
  description = "Test Item Type Description"
  modelKey = "${ox_model.Test_Model.key}"
}

resource "ox_link_type" "Test_Link_Type" {
  key = "test_link_type"
  name = "Test Link Type"
  description = "Test Link Type Description"
  modelKey = "${ox_model.Test_Model.key}"
}

resource "ox_link_rule" "Item_To_Item_Rule" {
  key = "Test_Item_Type->Test_Item_Type"
  name = "Test Item Type to Test Item Type rule"
  description = "Allow linking two items of type Test Item Type"
  linkTypeKey = "${ox_link_type.Test_Link_Type.key}"
  startItemTypeKey = "${ox_item_type.Test_Item_Type.key}"
  endItemTypeKey = "${ox_item_type.Test_Item_Type.key}"
}

resource "ox_item" "Item_1" {
  key = "item_1"
  name = "Item 1"
  description = "Item 1 Description"
  itemtype = "${ox_item_type.Test_Item_Type.key}"
  meta = <<JSON
    { "OS": "RHEL7.3" }
    JSON
}

resource "ox_item" "Item_2" {
  key = "item_2"
  name = "Item 2"
  description = "Item 2 Description"
  itemtype = "${ox_item_type.Test_Item_Type.key}"
  meta = <<JSON
    { "VM": true }
    JSON
}

resource "ox_link" "Link_1" {
  key = "link_1"
  description = "link 1 description"
  type = "${ox_link_type.Test_Link_Type.key}"
  startItemKey = "${ox_item.Item_1.key}"
  endItemKey = "${ox_item.Item_2.key}"
}