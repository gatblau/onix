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
	LinkDataSourceName  = "data.ox_link.ox_link_1_data"
	LinkDsItem1Key      = "test_acc_ox_link_1_item_1_data"
	LinkDsItem1Name     = "test_acc_ox_link_1_item_1_data name"
	LinkDsItem1Desc     = "test_acc_ox_link_1_item_1_data description"
	LinkDsItem2Key      = "test_acc_ox_link_1_item_2_data"
	LinkDsItem2Name     = "test_acc_ox_link_1_item_2_data name"
	LinkDsItem2Desc     = "test_acc_ox_link_1_item_2_data description"
	LinkDsKey           = "test_acc_ox_link_1_data"
	LinkDsDesc          = "ox_link_data Description"
	LinkDsMeta          = `{ "OS" = "RHEL7.3" }`
	LinkDsAttr          = `{ "RAM" : "3", "CPU" : "1" }`
	LinkDsLinkRuleKey   = "test_acc_item_type_ox_link_data-test_acc_item_type_ox_link_data"
	LinkDsLinkRuleName  = "test_acc_item_type_ox_link_data->test_acc_item_type_ox_link_data"
	LinkDsLinkRuleDesc  = "test_acc_item_type_ox_link_data->test_acc_item_type_ox_link_data description"
	LinkDsModelKey      = "test_acc_model_ox_link_data"
	LinkDsModelName     = "Model - ox_link_data"
	LinkDsModelDesc     = "Model for testing ox_link_data."
	LinkDsItemTypeKey   = "test_acc_item_type_ox_link_data"
	LinkDsItemTypeName  = "test_acc_item_type_ox_link_data name"
	LinkDsItemTypeDesc  = "test_acc_item_type_ox_link_data description"
	LinkDsLinkTypeKey   = "test_acc_link_type_ox_link_data"
	LinkDsLinkTypeName  = "test_acc_link_type_ox_link_data name"
	LinkDsLinkTypeDesc  = "test_acc_link_type_ox_link_data description"
	LinkDsTypeAttr1Key  = "test_acc_link_type_attr_cpu_ox_link_data"
	LinkDsTypeAttr1Name = "CPU"
	LinkDsTypeAttr1Desc = "CPU attr for testing ox_link."
	LinkDsTypeAttr2Key  = "test_acc_link_type_attr_ram_ox_link"
	LinkDsTypeAttr2Name = "RAM"
	LinkDsTypeAttr2Desc = "RAM attr for testing ox_link."
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("LinkDataSource", &resource.Sweeper{
		Name: "LinkDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteItem(&oxc.Item{Key: LinkDsItem1Key})
			_, err = cfg.Client.DeleteItem(&oxc.Item{Key: LinkDsItem2Key})
			_, err = cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkDsLinkTypeKey})
			_, err = cfg.Client.DeleteItemType(&oxc.ItemType{Key: LinkDsItemTypeKey})
			_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkDsModelKey})
			return err
		},
	})
}

func TestLinkDataSource(t *testing.T) {
	resourceName := LinkDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxLinkDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkDsKey),
					resource.TestCheckResourceAttr(resourceName, "description", LinkDsDesc),
					resource.TestCheckResourceAttr(resourceName, "type", LinkDsLinkTypeKey),
					resource.TestCheckResourceAttr(resourceName, "start_item_key", LinkDsItem1Key),
					resource.TestCheckResourceAttr(resourceName, "end_item_key", LinkDsItem2Key),
				),
			},
		},
	})
}

func oxLinkDataSource() string {
	return fmt.Sprintf(
		`data "ox_link" "ox_link_1_data" {
					key = "%s"	
				}`, LinkDsKey)
}

// create supporting database entities
func prepareLinkDataSourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkDsModelKey,
			Name:        LinkDsModelName,
			Description: LinkDsModelDesc,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// an item type in the model
	result, err = cfg.Client.PutItemType(
		&oxc.ItemType{
			Key:         LinkDsItemTypeKey,
			Name:        LinkDsItemTypeName,
			Description: LinkDsItemTypeDesc,
			Model:       LinkDsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// an link type in the model
	result, err = cfg.Client.PutLinkType(
		&oxc.LinkType{
			Key:         LinkDsLinkTypeKey,
			Name:        LinkDsLinkTypeName,
			Description: LinkDsLinkTypeDesc,
			Model:       LinkDsModelKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// and two attributes for the link type
	result, err = cfg.Client.PutLinkTypeAttr(
		&oxc.LinkTypeAttribute{
			Key:         LinkDsTypeAttr1Key,
			Name:        LinkDsTypeAttr1Name,
			Description: LinkDsTypeAttr1Desc,
			LinkTypeKey: LinkDsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	result, err = cfg.Client.PutLinkTypeAttr(
		&oxc.LinkTypeAttribute{
			Key:         LinkDsTypeAttr2Key,
			Name:        LinkDsTypeAttr2Name,
			Description: LinkDsTypeAttr2Desc,
			LinkTypeKey: LinkDsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// a link rule to allow to connect two item types
	result, err = cfg.Client.PutLinkRule(
		&oxc.LinkRule{
			Key:              LinkDsLinkRuleKey,
			Name:             LinkDsLinkRuleName,
			Description:      LinkDsLinkRuleDesc,
			StartItemTypeKey: LinkDsItemTypeKey,
			EndItemTypeKey:   LinkDsItemTypeKey,
			LinkTypeKey:      LinkDsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// item 1
	result, err = cfg.Client.PutItem(
		&oxc.Item{
			Key:         LinkDsItem1Key,
			Name:        LinkDsItem1Name,
			Description: LinkDsItem1Desc,
			Type:        LinkDsItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// item 2
	result, err = cfg.Client.PutItem(
		&oxc.Item{
			Key:         LinkDsItem2Key,
			Name:        LinkDsItem2Name,
			Description: LinkDsItem2Desc,
			Type:        LinkDsItemTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}

	// create the link between the two items
	result, err = cfg.Client.PutLink(
		&oxc.Link{
			Key:          LinkDsKey,
			Description:  LinkDsDesc,
			Type:         LinkDsLinkTypeKey,
			StartItemKey: LinkDsItem1Key,
			EndItemKey:   LinkDsItem2Key,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}
