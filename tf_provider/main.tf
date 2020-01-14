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
  //  secret    = "application-secret-here"
  //  token_uri = "uri-of-the-token-endpoint-at-authorisation-server"
}

resource "ox_model" "Test_Model" {
  key         = "test_model"
  name        = "Test Model"
  description = "Test Model Description"
}

resource "ox_item_type" "Test_Item_Type" {
  key         = "test_item_type"
  name        = "Test Item Type"
  description = "Test Item Type Description"
  model_key   = ox_model.Test_Model.key
}

resource "ox_link_type" "Test_Link_Type" {
  key         = "test_link_type"
  name        = "Test Link Type"
  description = "Test Link Type Description"
  model_key   = ox_model.Test_Model.key
}

resource "ox_link_rule" "Item_To_Item_Rule" {
  key                 = "Test_Item_Type->Test_Item_Type"
  name                = "Test Item Type to Test Item Type rule"
  description         = "Allow linking two items of type Test Item Type"
  link_type_key       = ox_link_type.Test_Link_Type.key
  start_item_type_key = ox_item_type.Test_Item_Type.key
  end_item_type_key   = ox_item_type.Test_Item_Type.key
}

resource "ox_item" "Item_1" {
  key         = "item_1"
  name        = "Item 1"
  description = "Item 1 Description"
  type        = ox_item_type.Test_Item_Type.key
  meta = {
    "OS" = "RHEL7.3"
  }

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

  // adds explicit dependency so that link rule is created first!
  depends_on = [ox_link_rule.Item_To_Item_Rule]
}

resource "ox_link" "Link_1" {
  key            = "link_1"
  description    = "link 1 description"
  type           = ox_link_type.Test_Link_Type.key
  start_item_key = ox_item.Item_1.key
  end_item_key   = ox_item.Item_2.key
}
