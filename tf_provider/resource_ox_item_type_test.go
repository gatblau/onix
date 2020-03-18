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
	"strconv"
	"testing"
)

// all the constants for the item test
const (
	ItemTypeResourceName   = "ox_item_type.ox_item_type_1"
	ItemTypeRsModelKey     = "test_acc_ox_item_type_model_1"
	ItemTypeRsModelName    = "ox_item_type_model_1"
	ItemTypeRsModelDesc    = "ox_item_type_model_1 Description"
	ItemTypeRsKey          = "test_acc_ox_item_type_1"
	ItemTypeRsName         = "ox_item_type_1_name"
	ItemTypeRsDesc         = "ox_item_type_1_description"
	ItemTypeRsNotifyChange = true
	ItemTypeRsEncryptMeta  = false
	ItemTypeRsEncryptTxt   = false
	ItemTypeRsManaged      = true
)

func TestItemTypeResource(t *testing.T) {
	resourceName := ItemTypeResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemTypeResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxItemTypeResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item type resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemTypeRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemTypeRsName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemTypeRsDesc),
					resource.TestCheckResourceAttr(resourceName, "model_key", ItemTypeRsModelKey),
					resource.TestCheckResourceAttr(resourceName, "notify_change", strconv.FormatBool(ItemTypeRsNotifyChange)),
					resource.TestCheckResourceAttr(resourceName, "encrypt_meta", strconv.FormatBool(ItemTypeRsEncryptMeta)),
					resource.TestCheckResourceAttr(resourceName, "encrypt_txt", strconv.FormatBool(ItemTypeRsEncryptTxt)),
					resource.TestCheckResourceAttr(resourceName, "managed", strconv.FormatBool(ItemTypeRsManaged)),

					// check for side effects in Onix database
					checkItemTypeResourceCreated(ItemTypeResourceName),
				),
			},
		},
		CheckDestroy: checkItemTypeResourceDestroyed,
	})
}

// create supporting database entities
func prepareItemTypeResourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemTypeRsModelKey,
			Name:        ItemTypeRsModelName,
			Description: ItemTypeRsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}

// check the item has been created in the database
func checkItemTypeResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the item from the database
		key := rs.Primary.Attributes["key"]
		itemType, err := cfg.Client.GetItemType(&oxc.ItemType{Key: key})
		if err != nil {
			return fmt.Errorf("can't read item type: %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "name", itemType.Name); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "description", itemType.Description); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "model_key", itemType.Model); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "managed", strconv.FormatBool(itemType.Managed)); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "encrypt_meta", strconv.FormatBool(itemType.EncryptMeta)); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "encrypt_txt", strconv.FormatBool(itemType.EncryptTxt)); err != nil {
			return err
		}
		if err := checkEntityAttr(rs, "notify_change", strconv.FormatBool(itemType.NotifyChange)); err != nil {
			return err
		}
		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkItemTypeResourceDestroyed(s *terraform.State) error {
	itemType, _ := cfg.Client.GetItemType(&oxc.ItemType{Key: ItemTypeRsKey})
	// should get an error as item should not exist anymore
	if itemType != nil {
		return fmt.Errorf("item yype %s still exists", ItemRsKey)
	}
	// clean up the database after the test
	return cleanUpItemTypeRs()
}

// remove supporting database entities
func cleanUpItemTypeRs() error {
	_, err := cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemTypeRsKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemTypeRsModelKey})
	return err
}

func oxItemTypeResource() string {
	return fmt.Sprintf(
		`resource "ox_item_type" "ox_item_type_1" {
  key         	= "%s"
  name        	= "%s"
  description 	= "%s"
  model_key   	= "%s"
  notify_change = %s
  encrypt_txt 	= %s
  encrypt_meta 	= "%s"
  managed 		= %s
}`, ItemTypeRsKey,
		ItemTypeRsName,
		ItemTypeRsDesc,
		ItemTypeRsModelKey,
		strconv.FormatBool(ItemTypeRsNotifyChange),
		strconv.FormatBool(ItemTypeRsEncryptTxt),
		strconv.FormatBool(ItemTypeRsEncryptMeta),
		strconv.FormatBool(ItemTypeRsManaged))
}
