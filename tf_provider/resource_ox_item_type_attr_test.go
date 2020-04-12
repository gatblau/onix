/*
   Onix Config Manager - Terraform Provider
   Copyright (c) 2018-2020 by www.gatblau.org

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software distributed under
   the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
   either express or implied.
   See the License for the specific language governing permissions and limitations under the License.

   Contributors to this project, hereby assign copyright in this code to the project,
   to be licensed under the same terms as the rest of the code.
*/
package main

import (
	"fmt"
	"github.com/gatblau/oxc"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

// all the constants for the item test
const (
	ItemTypeAttrResourceName   = "ox_item_type_attr.ox_item_type_attr_1"
	ItemTypeAttrRsKey          = "test_acc_ox_item_type_attr_1"
	ItemTypeAttrRsName         = "test_acc_ox_item_type_attr_1 name"
	ItemTypeAttrRsDesc         = "test_acc_ox_item_type_attr_1 description"
	ItemTypeAttrRsType         = "list"
	ItemTypeAttrRsDefValue     = "A,B,C"
	ItemTypeAttrRsItemTypeKey  = "test_acc_item_type_ox_item_type_attr"
	ItemTypeAttrRsItemTypeName = "test_acc_item_type_ox_item_type_attr name"
	ItemTypeAttrRsItemTypeDesc = "test_acc_item_type_ox_item_type_attr description"
	ItemTypeAttrRsModelKey     = "test_acc_model_ox_item_type_attr"
	ItemTypeAttrRsModelName    = "test_acc_model_ox_item_type_attr name"
	ItemTypeAttrRsModelDesc    = "test_acc_model_ox_item_type_attr description"
)

func TestItemTypeAttrResource(t *testing.T) {
	resourceName := ItemTypeAttrResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemTypeAttrResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxItemTypeAttrResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemTypeAttrRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemTypeAttrRsName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemTypeAttrRsDesc),
					resource.TestCheckResourceAttr(resourceName, "item_type_key", ItemTypeAttrRsItemTypeKey),
					resource.TestCheckResourceAttr(resourceName, "type", ItemTypeAttrRsType),
					resource.TestCheckResourceAttr(resourceName, "def_value", ItemTypeAttrRsDefValue),

					// check for side effects in Onix database
					checkItemTypeAttrResourceCreated(ItemTypeAttrResourceName),
				),
			},
		},
		CheckDestroy: checkItemTypeAttrResourceDestroyed,
	})
}

// create supporting database entities
func prepareItemTypeAttrResourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemTypeAttrRsModelKey,
			Name:        ItemTypeAttrRsModelName,
			Description: ItemTypeAttrRsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an item type in the model
	_, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         ItemTypeAttrRsItemTypeKey,
			Name:        ItemTypeAttrRsItemTypeName,
			Description: ItemTypeAttrRsItemTypeDesc,
			Model:       ItemTypeAttrRsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
}

// check the item has been created in the database
func checkItemTypeAttrResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the attr from the database
		key := rs.Primary.Attributes["key"]
		attr, err := cfg.Client.GetItemTypeAttr(&oxc.ItemTypeAttribute{Key: key, ItemTypeKey: ItemTypeAttrRsItemTypeKey})
		if err != nil {
			return fmt.Errorf("can't read attr type attribute %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "name", attr.Name); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "description", attr.Description); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "type", attr.Type); err != nil {
			return err
		}

		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkItemTypeAttrResourceDestroyed(s *terraform.State) error {
	item, _ := cfg.Client.GetItemTypeAttr(&oxc.ItemTypeAttribute{Key: ItemTypeAttrRsKey, ItemTypeKey: ItemTypeAttrRsItemTypeKey})
	// should get an error as item should not exist anymore
	if item != nil {
		return fmt.Errorf("item type attribute '%s' still exists", ItemTypeAttrRsKey)
	}
	// clean up the database after the test
	return cleanUpItemTypeAttrRs()
}

// remove supporting database entities
func cleanUpItemTypeAttrRs() error {
	_, err := cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemTypeAttrRsItemTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemTypeAttrRsModelKey})
	return err
}

func oxItemTypeAttrResource() string {
	return fmt.Sprintf(
		`resource "ox_item_type_attr" "ox_item_type_attr_1" {
  key         	= "%s"
  name        	= "%s"
  description 	= "%s"
  item_type_key = "%s"
  type        	= "%s"
  def_value   	= "%s"
}`, ItemTypeAttrRsKey,
		ItemTypeAttrRsName,
		ItemTypeAttrRsDesc,
		ItemTypeAttrRsItemTypeKey,
		ItemTypeAttrRsType,
		ItemTypeAttrRsDefValue)
}
