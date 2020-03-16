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
	ItemResourceName  = "ox_item.ox_item_1"
	ItemKey           = "ox_item_1"
	ItemName          = "ox item"
	ItemDesc          = "ox_item Description"
	ItemMeta          = `{ "OS" = "RHEL7.3" }`
	ItemAttr          = `{ "RAM" : "3", "CPU" : "1" }`
	ItemPartition     = "INS"
	ItemModelKey      = "model_ox_item"
	ItemModelName     = "Model - ox_item"
	ItemModelDesc     = "Model for testing ox_item."
	ItemTypeKey       = "item_type_ox_item"
	ItemTypeName      = "Item Type - ox_item"
	ItemTypeDesc      = "Item Type for testing ox_item."
	ItemTypeAttr1Key  = "item_type_attr_cpu_ox_item"
	ItemTypeAttr1Name = "CPU"
	ItemTypeAttr1Desc = "CPU attr for testing ox_item."
	ItemTypeAttr2Key  = "item_type_attr_ram_ox_item"
	ItemTypeAttr2Name = "RAM"
	ItemTypeAttr2Desc = "RAM attr for testing ox_item."
)

func TestItem(t *testing.T) {
	resourceName := ItemResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxItemResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemDesc),
					resource.TestCheckResourceAttr(resourceName, "type", ItemTypeKey),

					// check for side effects in Onix database
					checkItemCreated(ItemResourceName),
				),
			},
		},
		CheckDestroy: checkItemDestroyed,
	})
}

// create supporting database entities
func prepareItemTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemModelKey,
			Name:        ItemModelName,
			Description: ItemModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an item type in the model
	_, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         ItemTypeKey,
			Name:        ItemTypeName,
			Description: ItemTypeDesc,
			Model:       ItemModelKey,
		})
	if err != nil {
		t.Error(err)
	}

	// and two attributes for the item type
	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemTypeAttr1Key,
			Name:        ItemTypeAttr1Name,
			Description: ItemTypeAttr1Desc,
			ItemTypeKey: ItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}

	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemTypeAttr2Key,
			Name:        ItemTypeAttr2Name,
			Description: ItemTypeAttr2Desc,
			ItemTypeKey: ItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
}

// check the item has been created in the database
func checkItemCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		// retrieve the item from the database
		key := rs.Primary.Attributes["key"]
		item, err := cfg.Client.GetItem(&oxc.Item{Key: key})
		if err != nil {
			return fmt.Errorf("Can't read item %s", key)
		}

		// compares the state with the database values
		if err := checkItemAttr(rs, "name", item.Name); err != nil {
			return err
		}
		if err := checkItemAttr(rs, "description", item.Description); err != nil {
			return err
		}
		if err := checkItemAttr(rs, "partition", item.Partition); err != nil {
			return err
		}
		if err := checkItemAttr(rs, "type", item.Type); err != nil {
			return err
		}

		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkItemDestroyed(s *terraform.State) error {
	item, err := cfg.Client.GetItem(&oxc.Item{Key: ItemKey})
	// should get an error as item should not exist anymore
	if err == nil || len(item.Key) > 0 {
		return fmt.Errorf("Item %s still exists", ItemKey)
	}
	// clean up the database after the test
	return cleanUpItem()
}

// remove supporting database entities
func cleanUpItem() error {
	_, err := cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemModelKey})
	return err
}

// check the attribute in TF state matches the one in the database
func checkItemAttr(rs *terraform.ResourceState, attrName string, targetValue string) error {
	if rs.Primary.Attributes[attrName] != targetValue {
		return fmt.Errorf("Attribute '%s' expected value %s, but found %s", attrName, rs.Primary.Attributes[attrName], targetValue)
	}
	return nil
}

func oxItemResource() string {
	return fmt.Sprintf(
		`resource "ox_item" "ox_item_1" {
  key         = "%s"
  name        = "%s"
  description = "%s"
  type        = "%s"
  meta        = %s
  attribute   = %s
  partition   = "%s"
}`, ItemKey, ItemName, ItemDesc, ItemTypeKey, ItemMeta, ItemAttr, ItemPartition)
}
