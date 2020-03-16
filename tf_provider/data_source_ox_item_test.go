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

const (
	ItemDataSourceName  = "data.ox_item.ox_item_1_data"
	ItemDsKey           = "ox_item_1_data"
	ItemDsName          = "ox item 1 data"
	ItemDsDesc          = "ox_item_1_data Description"
	ItemDsText          = `free text for ox_item_1_data`
	ItemDsPartition     = "INS"
	ItemDsStatus        = 3
	ItemDsModelKey      = "model_ox_item_data"
	ItemDsModelName     = "Model - ox_item_data"
	ItemDsModelDesc     = "Model for testing ox_item_data."
	ItemDsTypeKey       = "item_type_ox_item_data"
	ItemDsTypeName      = "Item Type - ox_item_data"
	ItemDsTypeDesc      = "Item Type for testing ox_item_data."
	ItemDsTypeAttr1Key  = "item_type_attr_cpu_ox_item_data"
	ItemDsTypeAttr1Name = "CPU"
	ItemDsTypeAttr1Desc = "CPU attr for testing ox_item_data."
	ItemDsTypeAttr2Key  = "item_type_attr_ram_ox_item_data"
	ItemDsTypeAttr2Name = "RAM"
	ItemDsTypeAttr2Desc = "RAM attr for testing ox_item_data."
)

var (
	ItemDsTag  = []interface{}{"VM", "DC1", "IRELAND"}
	ItemDsMeta = make(map[string]interface{})
	ItemDsAttr = make(map[string]interface{})
)

func TestItemDataSource(t *testing.T) {
	resourceName := ItemDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxItemDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemDsName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemDsDesc),
					resource.TestCheckResourceAttr(resourceName, "type", ItemDsTypeKey),
					resource.TestCheckResourceAttr(resourceName, "status", strconv.Itoa(ItemDsStatus)),
					resource.TestCheckResourceAttr(resourceName, "partition", ItemDsPartition),
					resource.TestCheckResourceAttrSet(resourceName, "attribute.CPU"), // ItemDsTypeAttr1Name
					resource.TestCheckResourceAttrSet(resourceName, "attribute.RAM"), // ItemDsTypeAttr2Name
					resource.TestCheckResourceAttrSet(resourceName, "tag.0"),
					resource.TestCheckResourceAttrSet(resourceName, "tag.1"),
					resource.TestCheckResourceAttrSet(resourceName, "tag.2"),
					resource.TestCheckResourceAttrSet(resourceName, "meta.%"),
				),
			},
		},
	})
}

func oxItemDataSource() string {
	return fmt.Sprintf(
		`data "ox_item" "ox_item_1_data" {
					key = "%s"	
				}`, ItemDsKey)
}

// create supporting database entities
func prepareItemDataSourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemDsModelKey,
			Name:        ItemDsModelName,
			Description: ItemDsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an item type in the model
	_, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         ItemDsTypeKey,
			Name:        ItemDsTypeName,
			Description: ItemDsTypeDesc,
			Model:       ItemDsModelKey,
		})
	if err != nil {
		t.Error(err)
	}

	// and two attributes for the item type
	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemDsTypeAttr1Key,
			Name:        ItemDsTypeAttr1Name,
			Description: ItemDsTypeAttr1Desc,
			ItemTypeKey: ItemDsTypeKey,
		})
	if err != nil {
		t.Error(err)
	}

	_, err = cfg.Client.PutItemTypeAttr(
		&oxc.ItemTypeAttribute{
			Key:         ItemDsTypeAttr2Key,
			Name:        ItemDsTypeAttr2Name,
			Description: ItemDsTypeAttr2Desc,
			ItemTypeKey: ItemDsTypeKey,
		})
	if err != nil {
		t.Error(err)
	}

	// create item
	ItemDsMeta["OS"] = "RHEL 8"
	ItemDsAttr["RAM"] = 3
	ItemDsAttr["CPU"] = 2

	_, err = cfg.Client.PutItem(
		&oxc.Item{
			Key:         ItemDsKey,
			Name:        ItemDsName,
			Description: ItemDsDesc,
			Status:      ItemDsStatus,
			Type:        ItemDsTypeKey,
			Tag:         ItemDsTag,
			Meta:        ItemDsMeta,
			Txt:         ItemDsText,
			Attribute:   ItemDsAttr,
			Partition:   ItemDsPartition,
		})
	if err != nil {
		t.Error(err)
	}
}

// remove supporting database entities
func cleanUpItemDs(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := cfg.Client.DeleteItem(&oxc.Item{Key: ItemDsKey})
		_, err = cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemDsTypeKey})
		_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemDsModelKey})
		return err
	}
}
