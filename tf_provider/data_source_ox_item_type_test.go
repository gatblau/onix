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
	"testing"
)

const (
	ItemTypeDatasourceName = "data.ox_item_type.ox_item_type_1"
	ItemTypeDsModelKey     = "test_acc_ox_item_type_model_1_data"
	ItemTypeDsModelName    = "ox_item_type_model_1_data"
	ItemTypeDsModelDesc    = "ox_item_type_model_1_data Description"
	ItemTypeDsKey          = "test_acc_ox_item_type_1_data"
	ItemTypeDsName         = "ox_item_type_1_data_name"
	ItemTypeDsDesc         = "ox_item_type_1_data_description"
	ItemTypeDsNotifyChange = true
	ItemTypeDsEncryptMeta  = false
	ItemTypeDsEncryptTxt   = false
	ItemTypeDsManaged      = true
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("ItemTypeDataSource", &resource.Sweeper{
		Name: "ItemTypeDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteItemType(&oxc.ItemType{Key: ItemTypeDsKey})
			_, err = cfg.Client.DeleteModel(&oxc.Model{Key: ItemTypeDsModelKey})
			return err
		},
	})
}

func TestItemTypeDataSource(t *testing.T) {
	resourceName := ItemTypeDatasourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareItemTypeDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxItemTypeDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item type attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", ItemTypeDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", ItemTypeDsName),
					resource.TestCheckResourceAttr(resourceName, "description", ItemTypeDsDesc),
				),
			},
		},
	})
}

func oxItemTypeDataSource() string {
	return fmt.Sprintf(
		`data "ox_item_type" "ox_item_type_1" {
					key = "%s"	
				}`, ItemTypeDsKey)
}

// create supporting database entities
func prepareItemTypeDataSourceTest(t *testing.T) {
	// need a model first
	_, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         ItemTypeDsModelKey,
			Name:        ItemTypeDsModelName,
			Description: ItemTypeDsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}

	// an item type in the model
	_, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         ItemTypeDsKey,
			Name:        ItemTypeDsName,
			Description: ItemTypeDsDesc,
			Model:       ItemTypeDsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
}
