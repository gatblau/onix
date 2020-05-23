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
	LinkRuleDataSourceName = "data.ox_link_rule.ox_link_rule_1_data"
	LinkRuleDsKey          = "test_acc_item_type_ox_link_rule_data-test_acc_item_type_ox_link_rule_data"
	LinkRuleDsName         = "test_acc_item_type_ox_link_rule_data->test_acc_item_type_ox_link_rule_data name"
	LinkRuleDsDesc         = "test_acc_item_type_ox_link_rule_data->test_acc_item_type_ox_link_rule_data description"
	LinkRuleDsModelKey     = "test_acc_model_ox_link_rule_data"
	LinkRuleDsModelName    = "Model - ox_link_rule_data"
	LinkRuleDsModelDesc    = "Model for testing ox_link_rule_data."
	LinkRuleDsItemTypeKey  = "test_acc_item_type_ox_link_rule_data"
	LinkRuleDsItemTypeName = "test_acc_item_type_ox_link_rule_data name"
	LinkRuleDsItemTypeDesc = "test_acc_item_type_ox_link_rule_data description"
	LinkRuleDsLinkTypeKey  = "test_acc_link_type_ox_link_rule_data"
	LinkRuleDsLinkTypeName = "test_acc_link_type_ox_link_rule_data name"
	LinkRuleDsLinkTypeDesc = "test_acc_link_type_ox_link_rule_data description"
)

func init() {
	// defines a sweeper to clean up dangling test resources
	resource.AddTestSweepers("LinkRuleDataSource", &resource.Sweeper{
		Name: "LinkRuleDataSource",
		F: func(region string) error {
			_, err := cfg.Client.DeleteLinkRule(&oxc.LinkRule{Key: LinkRuleDsKey})
			_, err = cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkRuleDsLinkTypeKey})
			_, err = cfg.Client.DeleteItemType(&oxc.ItemType{Key: LinkRuleDsItemTypeKey})
			_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkRuleDsModelKey})
			return err
		},
	})
}

func TestLinkRuleDataSource(t *testing.T) {
	resourceName := LinkRuleDataSourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkRuleDataSourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// read
			{
				Config: oxLinkRuleDataSource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkRuleDsKey),
					resource.TestCheckResourceAttr(resourceName, "name", LinkRuleDsName),
					resource.TestCheckResourceAttr(resourceName, "description", LinkRuleDsDesc),
					resource.TestCheckResourceAttr(resourceName, "link_type_key", LinkRuleDsLinkTypeKey),
					resource.TestCheckResourceAttr(resourceName, "start_item_type_key", LinkRuleDsItemTypeKey),
					resource.TestCheckResourceAttr(resourceName, "end_item_type_key", LinkRuleDsItemTypeKey),
				),
			},
		},
	})
}

func oxLinkRuleDataSource() string {
	return fmt.Sprintf(
		`data "ox_link_rule" "ox_link_rule_1_data" {
					key = "%s"	
				}`, LinkRuleDsKey)
}

// create supporting database entities
func prepareLinkRuleDataSourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkRuleDsModelKey,
			Name:        LinkRuleDsModelName,
			Description: LinkRuleDsModelDesc,
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
			Key:         LinkRuleDsItemTypeKey,
			Name:        LinkRuleDsItemTypeName,
			Description: LinkRuleDsItemTypeDesc,
			Model:       LinkRuleDsModelKey,
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
			Key:         LinkRuleDsLinkTypeKey,
			Name:        LinkRuleDsLinkTypeName,
			Description: LinkRuleDsLinkTypeDesc,
			Model:       LinkRuleDsModelKey,
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
			Key:              LinkRuleDsKey,
			Name:             LinkRuleDsName,
			Description:      LinkRuleDsDesc,
			StartItemTypeKey: LinkRuleDsItemTypeKey,
			EndItemTypeKey:   LinkRuleDsItemTypeKey,
			LinkTypeKey:      LinkRuleDsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}
