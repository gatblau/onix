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

// all the constants for the link test
const (
	LinkRuleResourceName   = "ox_link_rule.ox_link_rule_1"
	LinkRuleRsKey          = "test_acc_item_type_ox_link_rule-test_acc_item_type_ox_link_rule"
	LinkRuleRsName         = "test_acc_item_type_ox_link_rule->test_acc_item_type_ox_link_rule"
	LinkRuleRsDesc         = "test_acc_item_type_ox_link_rule->test_acc_item_type_ox_link_rule description"
	LinkRuleRsModelKey     = "test_acc_model_ox_link_rule"
	LinkRuleRsModelName    = "test_acc_model_ox_link_rule Name"
	LinkRuleRsModelDesc    = "test_acc_model_ox_link_rule Description"
	LinkRuleRsItemTypeKey  = "test_acc_item_type_ox_link_rule"
	LinkRuleRsItemTypeName = "test_acc_item_type_ox_link_rule name"
	LinkRuleRsItemTypeDesc = "test_acc_item_type_ox_link_rule description"
	LinkRuleRsLinkTypeKey  = "test_acc_link_type_ox_link_rule"
	LinkRuleRsLinkTypeName = "test_acc_link_type_ox_link_rule name"
	LinkRuleRsLinkTypeDesc = "test_acc_link_type_ox_link_rule description"
)

func TestLinkRuleResource(t *testing.T) {
	resourceName := LinkRuleResourceName
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { prepareLinkRuleResourceTest(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			//create
			{
				Config: oxLinkRuleResource(),
				Check: resource.ComposeTestCheckFunc(
					// check item resource attributes in Terraform state
					resource.TestCheckResourceAttr(resourceName, "key", LinkRuleRsKey),
					resource.TestCheckResourceAttr(resourceName, "name", LinkRuleRsName),
					resource.TestCheckResourceAttr(resourceName, "description", LinkRuleRsDesc),
					resource.TestCheckResourceAttr(resourceName, "link_type_key", LinkRuleRsLinkTypeKey),
					resource.TestCheckResourceAttr(resourceName, "start_item_type_key", LinkRuleRsItemTypeKey),
					resource.TestCheckResourceAttr(resourceName, "end_item_type_key", LinkRuleRsItemTypeKey),

					// check for side effects in Onix database
					checkLinkRuleResourceCreated(LinkRuleResourceName),
				),
			},
		},
		CheckDestroy: checkLinkRuleResourceDestroyed,
	})
}

// create supporting database entities
func prepareLinkRuleResourceTest(t *testing.T) {
	// need a model first
	result, err := cfg.Client.PutModel(
		&oxc.Model{
			Key:         LinkRuleRsModelKey,
			Name:        LinkRuleRsModelName,
			Description: LinkRuleRsModelDesc,
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
			Key:         LinkRuleRsItemTypeKey,
			Name:        LinkRuleRsItemTypeName,
			Description: LinkRuleRsItemTypeDesc,
			Model:       LinkRuleRsModelKey,
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
			Key:         LinkRuleRsLinkTypeKey,
			Name:        LinkRuleRsLinkTypeName,
			Description: LinkRuleRsLinkTypeDesc,
			Model:       LinkRuleRsModelKey,
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
			Key:              LinkRuleRsKey,
			Name:             LinkRuleRsName,
			Description:      LinkRuleRsDesc,
			StartItemTypeKey: LinkRuleRsItemTypeKey,
			EndItemTypeKey:   LinkRuleRsItemTypeKey,
			LinkTypeKey:      LinkRuleRsLinkTypeKey,
		})
	if err != nil {
		t.Error(err)
	}
	if result.Error {
		t.Error(result.Message)
	}
}

// check the link has been created in the database
func checkLinkRuleResourceCreated(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// retrieve the resource by name from state
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		// retrieve the link from the database
		key := rs.Primary.Attributes["key"]
		rule, err := cfg.Client.GetLinkRule(&oxc.LinkRule{Key: key})
		if err != nil {
			return fmt.Errorf("can't read link rule %s", key)
		}

		// compares the state with the database values
		if err := checkEntityAttr(rs, "description", rule.Description); err != nil {
			return err
		}
		return nil
	}
}

// check the item has been destroyed destroyed in the database
func checkLinkRuleResourceDestroyed(s *terraform.State) error {
	rule, _ := cfg.Client.GetLinkRule(&oxc.LinkRule{Key: LinkRuleRsKey})
	// should get an error as link rule should not exist anymore
	if rule != nil {
		return fmt.Errorf("rule rule %s still exists", LinkRuleRsKey)
	}
	// clean up the database after the test
	return cleanUpLinkRuleRs()
}

// remove supporting database entities
func cleanUpLinkRuleRs() error {
	_, err := cfg.Client.DeleteLinkType(&oxc.LinkType{Key: LinkRuleRsLinkTypeKey})
	_, err = cfg.Client.DeleteItemType(&oxc.ItemType{Key: LinkRuleRsItemTypeKey})
	_, err = cfg.Client.DeleteModel(&oxc.Model{Key: LinkRuleRsLinkTypeKey})
	return err
}

func oxLinkRuleResource() string {
	return fmt.Sprintf(`resource "ox_link_rule" "ox_link_rule_1" {
  key = "%s"
  name = "%s"
  description = "%s"
  link_type_key = "%s"
  start_item_type_key = "%s"
  end_item_type_key = "%s"
}`, LinkRuleRsKey,
		LinkRuleRsName,
		LinkRuleRsDesc,
		LinkRuleRsLinkTypeKey,
		LinkRuleRsItemTypeKey,
		LinkRuleRsItemTypeKey)
}
