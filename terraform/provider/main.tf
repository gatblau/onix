terraform {
  backend "http" {
    address = "http://localhost:8081/state/foo"
    lock_address = "http://localhost:8081/state/foo"
    unlock_address = "http://localhost:8081/state/foo"
    username = "admin"
    password = "0n1x"
  }
}

provider "ox" {
  uri = "http://localhost:8080"

  //  use basic authentication
  auth_mode = "basic"
  user      = "admin"
  pwd       = "0n1x"

  //  uncomment & reset values below to
  //  use OpenId authentication
  //  auth_mode = "oidc"
  //  client_id = "application-client-id-here"
  //  app_secret = "application-secret-here"
  //  token_uri = "uri-of-the-token-endpoint-at-authorisation-server"
}

resource "ox_partition" "Test_Partition" {
  key         = "TEST"
  name        = "Test Partition"
  description = "A logical partition for testing."
}

resource "ox_role" "Test_Role" {
  key         = "ADMIN"
  name        = "System Administrator"
  description = "Can read and write configuration data and models."
  level       = 2

  depends_on = [ox_partition.Test_Partition]
}

resource "ox_privilege" "Test_Privilege" {
  key        = "ADMIN-TEST"
  role       = ox_role.Test_Role.key
  partition  = ox_partition.Test_Partition.key
  can_create = true
  can_read   = true
  can_delete = true

  depends_on = [ox_role.Test_Role]
}

resource "ox_model" "Test_Model" {
  key         = "test_model"
  name        = "Test Model"
  description = "Test Model Description"
  managed     = false

  depends_on = [ox_privilege.Test_Privilege]
}

resource "ox_item_type" "Test_Item_Type" {
  key           = "test_item_type"
  name          = "Test Item Type"
  description   = "Test Item Type Description"
  model_key     = ox_model.Test_Model.key
  notify_change = true
  encrypt_txt   = true
  encrypt_meta  = false
  managed       = true

  depends_on = [ox_model.Test_Model]
}

resource "ox_item_type_attr" "Test_Item_Type_Attr_1" {
  key           = "test_item_type_attr_1"
  name          = "RAM"
  description   = "GB of RAM"
  item_type_key = ox_item_type.Test_Item_Type.key
  type          = "integer"
  def_value     = "2"
  managed       = false

  depends_on = [ox_item_type.Test_Item_Type]
}

resource "ox_item_type_attr" "Test_Item_Type_Attr_2" {
  key           = "test_item_type_attr_2"
  name          = "CPU"
  description   = "No of CPU"
  item_type_key = ox_item_type.Test_Item_Type.key
  type          = "integer"
  def_value     = "1"
  managed       = false

  depends_on = [ox_item_type_attr.Test_Item_Type_Attr_1]
}

resource "ox_link_type" "Test_Link_Type" {
  key          = "test_link_type"
  name         = "Test Link Type"
  description  = "Test Link Type Description"
  model_key    = ox_model.Test_Model.key
  encrypt_meta = false
  encrypt_txt  = false
  managed      = false
  depends_on   = [ox_item_type_attr.Test_Item_Type_Attr_2]
}

resource "ox_link_type_attr" "Test_Link_Type_Attr_1" {
  key           = "test_link_type_attr_1"
  name          = "TEAM"
  description   = "Team Name"
  type          = "string"
  def_value     = "A-Team"
  managed       = false
  link_type_key = ox_link_type.Test_Link_Type.key

  depends_on = [ox_link_type.Test_Link_Type]
}

resource "ox_link_type_attr" "Test_Link_Type_Attr_2" {
  key           = "test_link_type_attr_2"
  name          = "CATEGORY"
  description   = "Name of Category"
  type          = "string"
  link_type_key = ox_link_type.Test_Link_Type.key

  depends_on = [ox_link_type_attr.Test_Link_Type_Attr_1]
}

resource "ox_link_rule" "Item_To_Item_Rule" {
  key                 = "Test_Item_Type->Test_Item_Type"
  name                = "Test Item Type to Test Item Type rule"
  description         = "Allow linking two items of type Test Item Type"
  link_type_key       = ox_link_type.Test_Link_Type.key
  start_item_type_key = ox_item_type.Test_Item_Type.key
  end_item_type_key   = ox_item_type.Test_Item_Type.key

  depends_on = [ox_link_type_attr.Test_Link_Type_Attr_2]
}

resource "ox_item" "Item_1" {
  key         = "item_1"
  name        = "Item 1"
  description = "Item 1 Description"
  type        = ox_item_type.Test_Item_Type.key
  meta = {
    "OS" = "RHEL7.3"
  }
  attribute = {
    "RAM" : "3",
    "CPU" : "1"
  }
  partition = ox_partition.Test_Partition.key

  // adds explicit dependency so that link rule is created first!
  depends_on = [ox_link_rule.Item_To_Item_Rule]
}

resource "ox_item" "Item_2" {
  key         = "item_2"
  name        = "Item 2"
  description = "Item 2 Description"
  type        = ox_item_type.Test_Item_Type.key
  meta = {
    "VM" = true
  }
  attribute = {
    "RAM" : "3",
    "CPU" : "1"
  }
  partition = ox_partition.Test_Partition.key

  // adds explicit dependency so that link rule is created first!
  depends_on = [ox_link_rule.Item_To_Item_Rule]
}

resource "ox_link" "Link_1" {
  key            = "link_1"
  description    = "link 1 description"
  type           = ox_link_type.Test_Link_Type.key
  start_item_key = ox_item.Item_1.key
  end_item_key   = ox_item.Item_2.key

  attribute = {
    "TEAM" : "Blue",
    "CATEGORY" : "Social"
  }

  depends_on = [ox_link_type.Test_Link_Type]
}

data "ox_model" "test_model_data" {
  key = ox_model.Test_Model.key
}

data "ox_item_type" "test_item_type_data" {
  key = ox_item_type.Test_Item_Type.key
}

data "ox_link_type" "test_item_type_data" {
  key = ox_link_type.Test_Link_Type.key
}

data "ox_item" "item_1_data" {
  key = ox_item.Item_1.key
}

data "ox_item" "item_2_data" {
  key = ox_item.Item_2.key
}

data "ox_item_type_attr" "Test_Item_Type_Attr_1_data" {
  key           = ox_item_type_attr.Test_Item_Type_Attr_1.key
  item_type_key = ox_item_type.Test_Item_Type.key
}

data "ox_link_type_attr" "Test_Link_Type_Attr_1_data" {
  key           = ox_link_type_attr.Test_Link_Type_Attr_1.key
  link_type_key = ox_link_type.Test_Link_Type.key
}

data "ox_link_rule" "Item_To_Item_Rule_Data" {
  key = ox_link_rule.Item_To_Item_Rule.key
}

data "ox_item" "Item_1_Data" {
  key = ox_item.Item_1.key
}

data "ox_link" "Link_1_Data" {
  key = ox_link.Link_1.key
}