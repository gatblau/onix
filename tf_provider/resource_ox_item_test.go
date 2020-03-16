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
	ItemResourceName    = "ox_item.ox_item_1"
	ItemRsKey           = "ox_item_1"
	ItemRsName          = "ox item"
	ItemRsDesc          = "ox_item Description"
	ItemRsMeta          = `{ "OS" = "RHEL7.3" }`
	ItemRsAttr          = `{ "RAM" : "3", "CPU" : "1" }`
	ItemRsPartition     = "INS"
	ItemRsModelKey      = "model_ox_item"
	ItemRsModelName     = "Model - ox_item"
	ItemRsModelDesc     = "Model for testing ox_item."
	ItemRsTypeKey       = "item_type_ox_item"
	ItemRsTypeName      = "Item Type - ox_item"
	ItemRsTypeDesc      = "Item Type for testing ox_item."
	ItemRsTypeAttr1Key  = "item_type_attr_cpu_ox_item"
	ItemRsTypeAttr1Name = "CPU"
	ItemRsTypeAttr1Desc = "CPU attr for testing ox_item."
	ItemRsTypeAttr2Key  = "item_type_attr_ram_ox_item"
	ItemRsTypeAttr2Name = "RAM"
	ItemRsTypeAttr2Desc = "RAM attr for testing ox_item."
)

func TestItemResource(t *testing.T) {
	resourceName := ItemResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxItemResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemRsName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemRsDesc),
					resource.TestCheckResourceAttr(resourceName, "type", ItemRsTypeKey),

					// check for side effects in Onix database
					checkItemResourceCreated(ItemResourceName),
				),
			},
		},
		CheckDestroy: checkItemResourceDestroyed,
	})
}

// create supporting database entities
func prepareItemResourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemRsModelKey,
			Name:        ItemRsModelName,
			Description: ItemRsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an item type in the model
	_, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         ItemRsTypeKey,
			Name:        ItemRsTypeName,
			Description: ItemRsTypeDesc,
			Model:       ItemRsModelKey,
		})
	if err != nil {
		t.Error(err)
	}

	// and two attributes for the item type
	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemRsTypeAttr1Key,
			Name:        ItemRsTypeAttr1Name,
			Description: ItemRsTypeAttr1Desc,
			ItemTypeKey: ItemRsTypeKey,
		})
	if err != nil {
		t.Error(err)
	}

	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemRsTypeAttr2Key,
			Name:        ItemRsTypeAttr2Name,
			Description: ItemRsTypeAttr2Desc,
			ItemTypeKey: ItemRsTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
}

// check the item has been created in the database
func checkItemResourceCreated(resourceName string) resource.TestCheckFunc {
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
func checkItemResourceDestroyed(s *terraform.State) error {
	item, err := cfg.Client.GetItem(&oxc.Item{Key: ItemRsKey})
	// should get an error as item should not exist anymore
	if err == nil || len(item.Key) > 0 {
		return fmt.Errorf("Item %s still exists", ItemRsKey)
	}
	// clean up the database after the test
	return cleanUpItemRs()
}

// remove supporting database entities
func cleanUpItemRs() error {
	_, err := cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemRsTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemRsModelKey})
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
}`, ItemRsKey, ItemRsName, ItemRsDesc, ItemRsTypeKey, ItemRsMeta, ItemRsAttr, ItemRsPartition)
}
